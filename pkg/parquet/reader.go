package parquet

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/parquet-go/parquet-go"
)

// ParquetSchema represents the schema of a Parquet file
type ParquetSchema struct {
	Columns []ParquetColumn `json:"columns"`
}

// ParquetColumn represents a column in the schema
type ParquetColumn struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

// GetSchema extracts schema from Parquet data
func GetSchema(fileData []byte) (*ParquetSchema, error) {
	reader := bytes.NewReader(fileData)
	pf, err := parquet.OpenFile(reader, int64(len(fileData)))
	if err != nil {
		return nil, fmt.Errorf("failed to open parquet file: %w", err)
	}

	schema := &ParquetSchema{
		Columns: make([]ParquetColumn, 0),
	}

	// Get schema from the file
	rootSchema := pf.Schema()
	for _, field := range rootSchema.Fields() {
		col := ParquetColumn{
			Name:     field.Name(),
			Type:     getTypeName(field),
			Nullable: field.Optional(),
		}
		schema.Columns = append(schema.Columns, col)
	}

	return schema, nil
}

func getTypeName(field parquet.Field) string {
	if field.Leaf() {
		node := field.Type()
		if node != nil {
			return node.String()
		}
	}
	return "UNKNOWN"
}

// ReadParquetFromS3 reads a Parquet file from S3 and returns a data frame
func ReadParquetFromS3(
	ctx context.Context,
	s3Client *s3.Client,
	bucket, key string,
	timeColumn string,
	columns []string,
	filter string,
	maxRows int,
	from, to time.Time,
) (*data.Frame, error) {
	// Download file from S3
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer func() { _ = result.Body.Close() }()

	fileData, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return ReadParquetData(fileData, timeColumn, columns, filter, maxRows, from, to)
}

// ReadParquetData reads Parquet data from bytes
func ReadParquetData(
	fileData []byte,
	timeColumn string,
	selectedColumns []string,
	filter string,
	maxRows int,
	from, to time.Time,
) (*data.Frame, error) {
	reader := bytes.NewReader(fileData)
	pf, err := parquet.OpenFile(reader, int64(len(fileData)))
	if err != nil {
		return nil, fmt.Errorf("failed to open parquet file: %w", err)
	}

	// Get schema info
	schema := pf.Schema()
	fields := schema.Fields()

	// Determine which columns to read
	var columnNames []string
	if len(selectedColumns) > 0 {
		columnNames = selectedColumns
	} else {
		for _, field := range fields {
			columnNames = append(columnNames, field.Name())
		}
	}

	// Create column indices map
	columnIndices := make(map[string]int)
	for i, field := range fields {
		columnIndices[field.Name()] = i
	}

	// Create frame
	frame := data.NewFrame("result")

	// Initialize fields based on schema
	fieldMap := make(map[string]*data.Field)
	for _, colName := range columnNames {
		idx, ok := columnIndices[colName]
		if !ok {
			continue
		}
		field := fields[idx]
		dataField := createFieldForParquetType(colName, field)
		frame.Fields = append(frame.Fields, dataField)
		fieldMap[colName] = dataField
	}

	// Read rows
	rowCount := 0
	for _, rowGroup := range pf.RowGroups() {
		rows := rowGroup.Rows()
		defer func() { _ = rows.Close() }()

		for maxRows <= 0 || rowCount < maxRows {

			row := make([]parquet.Value, len(fields))
			n, err := rows.ReadRows([]parquet.Row{row})
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, fmt.Errorf("failed to read row: %w", err)
			}
			if n == 0 {
				break
			}

			// Convert row to map for filtering
			rowMap := make(map[string]interface{})
			for i, field := range fields {
				rowMap[field.Name()] = parquetValueToInterface(row[i])
			}

			// Apply time filter if we have time column
			if timeColumn != "" && !from.IsZero() && !to.IsZero() {
				if timeVal, ok := rowMap[timeColumn]; ok {
					t := toTime(timeVal)
					if t.Before(from) || t.After(to) {
						continue
					}
				}
			}

			// Apply custom filter
			if filter != "" && !matchesFilter(rowMap, filter) {
				continue
			}

			// Append values to fields
			for _, colName := range columnNames {
				if dataField, ok := fieldMap[colName]; ok {
					if val, ok := rowMap[colName]; ok {
						appendValue(dataField, val)
					}
				}
			}

			rowCount++
		}

		if maxRows > 0 && rowCount >= maxRows {
			break
		}
	}

	return frame, nil
}

func parquetValueToInterface(v parquet.Value) interface{} {
	if v.IsNull() {
		return nil
	}

	switch v.Kind() {
	case parquet.Boolean:
		return v.Boolean()
	case parquet.Int32:
		return v.Int32()
	case parquet.Int64:
		return v.Int64()
	case parquet.Float:
		return v.Float()
	case parquet.Double:
		return v.Double()
	case parquet.ByteArray, parquet.FixedLenByteArray:
		return string(v.ByteArray())
	default:
		return fmt.Sprintf("%v", v)
	}
}

func createFieldForParquetType(name string, field parquet.Field) *data.Field {
	if !field.Leaf() {
		return data.NewField(name, nil, []string{})
	}

	node := field.Type()
	if node == nil {
		return data.NewField(name, nil, []string{})
	}

	switch node.Kind() {
	case parquet.Boolean:
		return data.NewField(name, nil, []bool{})
	case parquet.Int32:
		return data.NewField(name, nil, []int32{})
	case parquet.Int64:
		return data.NewField(name, nil, []int64{})
	case parquet.Float:
		return data.NewField(name, nil, []float32{})
	case parquet.Double:
		return data.NewField(name, nil, []float64{})
	default:
		return data.NewField(name, nil, []string{})
	}
}

func appendValue(field *data.Field, val interface{}) {
	if val == nil {
		// Handle nil values by appending zero value
		switch field.Type() {
		case data.FieldTypeInt32:
			field.Append(int32(0))
		case data.FieldTypeInt64:
			field.Append(int64(0))
		case data.FieldTypeFloat32:
			field.Append(float32(0))
		case data.FieldTypeFloat64:
			field.Append(float64(0))
		case data.FieldTypeBool:
			field.Append(false)
		default:
			field.Append("")
		}
		return
	}

	switch field.Type() {
	case data.FieldTypeInt32:
		switch v := val.(type) {
		case int32:
			field.Append(v)
		case int64:
			field.Append(int32(v))
		case int:
			field.Append(int32(v))
		default:
			field.Append(int32(0))
		}
	case data.FieldTypeInt64:
		switch v := val.(type) {
		case int64:
			field.Append(v)
		case int32:
			field.Append(int64(v))
		case int:
			field.Append(int64(v))
		default:
			field.Append(int64(0))
		}
	case data.FieldTypeFloat32:
		switch v := val.(type) {
		case float32:
			field.Append(v)
		case float64:
			field.Append(float32(v))
		default:
			field.Append(float32(0))
		}
	case data.FieldTypeFloat64:
		switch v := val.(type) {
		case float64:
			field.Append(v)
		case float32:
			field.Append(float64(v))
		default:
			field.Append(float64(0))
		}
	case data.FieldTypeBool:
		if b, ok := val.(bool); ok {
			field.Append(b)
		} else {
			field.Append(false)
		}
	case data.FieldTypeTime:
		field.Append(toTime(val))
	default:
		field.Append(fmt.Sprintf("%v", val))
	}
}

func toTime(val interface{}) time.Time {
	switch v := val.(type) {
	case time.Time:
		return v
	case int64:
		// Parquet stores timestamps as nanoseconds
		return time.Unix(0, v)
	case int32:
		// Smaller values might be seconds or days
		return time.Unix(int64(v), 0)
	case float64:
		// Assume nanoseconds
		return time.Unix(0, int64(v))
	default:
		return time.Time{}
	}
}

func matchesFilter(row map[string]interface{}, filter string) bool {
	// Basic filter implementation
	// Supports: column = value, column > value, column < value
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return true
	}

	// Split by AND
	conditions := strings.Split(filter, " AND ")
	for _, cond := range conditions {
		cond = strings.TrimSpace(cond)
		if !evaluateCondition(row, cond) {
			return false
		}
	}
	return true
}

func evaluateCondition(row map[string]interface{}, cond string) bool {
	// Parse condition
	var col, op, value string

	for _, operator := range []string{">=", "<=", "!=", "=", ">", "<"} {
		parts := strings.SplitN(cond, operator, 2)
		if len(parts) == 2 {
			col = strings.TrimSpace(parts[0])
			op = operator
			value = strings.Trim(strings.TrimSpace(parts[1]), "'\"")
			break
		}
	}

	if col == "" {
		return true
	}

	rowVal, ok := row[col]
	if !ok {
		return false
	}

	// Compare based on type
	switch v := rowVal.(type) {
	case string:
		switch op {
		case "=":
			return v == value
		case "!=":
			return v != value
		}
	case int64, int32, int:
		var numVal int64
		_, _ = fmt.Sscanf(value, "%d", &numVal)
		rowNum := reflect.ValueOf(v).Int()
		return compareInt(rowNum, numVal, op)
	case float64, float32:
		var numVal float64
		_, _ = fmt.Sscanf(value, "%f", &numVal)
		rowNum := reflect.ValueOf(v).Float()
		return compareFloat(rowNum, numVal, op)
	}

	return true
}

func compareInt(a, b int64, op string) bool {
	switch op {
	case "=":
		return a == b
	case "!=":
		return a != b
	case ">":
		return a > b
	case "<":
		return a < b
	case ">=":
		return a >= b
	case "<=":
		return a <= b
	}
	return false
}

func compareFloat(a, b float64, op string) bool {
	switch op {
	case "=":
		return a == b
	case "!=":
		return a != b
	case ">":
		return a > b
	case "<":
		return a < b
	case ">=":
		return a >= b
	case "<=":
		return a <= b
	}
	return false
}
