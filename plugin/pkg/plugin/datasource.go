package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"

	parquetS3 "github.com/tobiasworkstech/grafana-parquet-s3-datasource/pkg/s3"
	parquetReader "github.com/tobiasworkstech/grafana-parquet-s3-datasource/pkg/parquet"
)

// Datasource is the main plugin datasource struct
type Datasource struct {
	settings DatasourceSettings
	s3Client *s3.Client
}

// DatasourceSettings contains the datasource configuration
type DatasourceSettings struct {
	Endpoint           string `json:"endpoint"`
	Region             string `json:"region"`
	AccessKeyID        string `json:"accessKeyId"`
	DefaultBucket      string `json:"defaultBucket"`
	PathStyle          bool   `json:"pathStyle"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
}

// SecureSettings contains encrypted settings
type SecureSettings struct {
	SecretAccessKey string `json:"secretAccessKey"`
}

// QueryModel represents a query from the frontend
type QueryModel struct {
	RefID       string   `json:"refId"`
	Bucket      string   `json:"bucket"`
	PathPrefix  string   `json:"pathPrefix"`
	FilePattern string   `json:"filePattern"`
	QueryText   string   `json:"queryText"`
	TimeColumn  string   `json:"timeColumn"`
	Columns     []string `json:"columns"`
	MaxRows     int      `json:"maxRows"`
	TimeRange   struct {
		From int64 `json:"from"`
		To   int64 `json:"to"`
	} `json:"timeRange"`
}

// SchemaRequest represents a schema request
type SchemaRequest struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

// FilesRequest represents a file listing request
type FilesRequest struct {
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
}

// NewDatasource creates a new datasource instance
func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var dsSettings DatasourceSettings
	if err := json.Unmarshal(settings.JSONData, &dsSettings); err != nil {
		return nil, fmt.Errorf("failed to parse settings: %w", err)
	}

	// Get secure settings
	secretAccessKey := settings.DecryptedSecureJSONData["secretAccessKey"]

	// Create S3 client
	s3Client := parquetS3.NewS3Client(
		dsSettings.Endpoint,
		dsSettings.Region,
		dsSettings.AccessKeyID,
		secretAccessKey,
		dsSettings.PathStyle,
		dsSettings.InsecureSkipVerify,
	)

	return &Datasource{
		settings: dsSettings,
		s3Client: s3Client,
	}, nil
}

// Dispose cleans up resources
func (d *Datasource) Dispose() {}

// QueryData handles multiple queries
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := d.query(ctx, req.PluginContext, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *Datasource) query(ctx context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	var response backend.DataResponse

	var qm QueryModel
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("failed to parse query: %v", err))
	}

	// Use default bucket if not specified
	bucket := qm.Bucket
	if bucket == "" {
		bucket = d.settings.DefaultBucket
	}
	if bucket == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "no bucket specified")
	}

	// Set defaults
	if qm.MaxRows <= 0 {
		qm.MaxRows = 10000
	}
	if qm.FilePattern == "" {
		qm.FilePattern = "*.parquet"
	}
	// TimeColumn can be empty for non-time-series data

	// List matching files
	files, err := d.listMatchingFiles(ctx, bucket, qm.PathPrefix, qm.FilePattern)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to list files: %v", err))
	}

	if len(files) == 0 {
		// Return empty frame
		frame := data.NewFrame(qm.RefID)
		response.Frames = append(response.Frames, frame)
		return response
	}

	// Determine time range - use zero time if not specified
	var fromTime, toTime time.Time
	if qm.TimeRange.From > 0 {
		fromTime = time.UnixMilli(qm.TimeRange.From)
	}
	if qm.TimeRange.To > 0 {
		toTime = time.UnixMilli(qm.TimeRange.To)
	}

	// Read and combine parquet files
	frame, err := d.readParquetFiles(ctx, bucket, files, qm, fromTime, toTime)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("failed to read parquet files: %v", err))
	}

	frame.RefID = qm.RefID
	response.Frames = append(response.Frames, frame)
	return response
}

func (d *Datasource) listMatchingFiles(ctx context.Context, bucket, prefix, pattern string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}
	if prefix != "" {
		input.Prefix = aws.String(prefix)
	}

	var files []string
	paginator := s3.NewListObjectsV2Paginator(d.s3Client, input)

	// Convert pattern to regex
	regexPattern := strings.ReplaceAll(pattern, ".", "\\.")
	regexPattern = strings.ReplaceAll(regexPattern, "*", ".*")
	regexPattern = "^" + regexPattern + "$"
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			filename := filepath.Base(key)
			if re.MatchString(filename) {
				files = append(files, key)
			}
		}
	}

	return files, nil
}

func (d *Datasource) readParquetFiles(ctx context.Context, bucket string, files []string, qm QueryModel, fromTime, toTime time.Time) (*data.Frame, error) {
	var allFrames []*data.Frame

	for _, file := range files {
		frame, err := parquetReader.ReadParquetFromS3(
			ctx,
			d.s3Client,
			bucket,
			file,
			qm.TimeColumn,
			qm.Columns,
			qm.QueryText,
			qm.MaxRows,
			fromTime,
			toTime,
		)
		if err != nil {
			log.DefaultLogger.Warn("Failed to read file", "file", file, "error", err)
			continue
		}
		allFrames = append(allFrames, frame)
	}

	if len(allFrames) == 0 {
		return data.NewFrame("result"), nil
	}

	// Combine frames
	return combineFrames(allFrames), nil
}

func combineFrames(frames []*data.Frame) *data.Frame {
	if len(frames) == 0 {
		return data.NewFrame("result")
	}
	if len(frames) == 1 {
		return frames[0]
	}

	// Use first frame as base
	result := frames[0]

	// Append data from other frames
	for i := 1; i < len(frames); i++ {
		for fieldIdx := range result.Fields {
			if fieldIdx < len(frames[i].Fields) {
				for rowIdx := 0; rowIdx < frames[i].Fields[fieldIdx].Len(); rowIdx++ {
					result.Fields[fieldIdx].Append(frames[i].Fields[fieldIdx].At(rowIdx))
				}
			}
		}
	}

	return result
}

// CheckHealth performs a health check
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	bucket := d.settings.DefaultBucket

	// Try to list buckets or access default bucket
	if bucket != "" {
		_, err := d.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return &backend.CheckHealthResult{
				Status:  backend.HealthStatusError,
				Message: fmt.Sprintf("Failed to access bucket '%s': %v", bucket, err),
			}, nil
		}
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusOk,
			Message: fmt.Sprintf("Successfully connected to S3 and verified access to bucket '%s'", bucket),
		}, nil
	}

	// Try to list buckets
	_, err := d.s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("Failed to connect to S3: %v", err),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Successfully connected to S3",
	}, nil
}

// CallResource handles resource calls from frontend
func (d *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	switch req.Path {
	case "health":
		return d.handleHealth(ctx, req, sender)
	case "buckets":
		return d.handleBuckets(ctx, req, sender)
	case "files":
		return d.handleFiles(ctx, req, sender)
	case "schema":
		return d.handleSchema(ctx, req, sender)
	case "query":
		return d.handleQuery(ctx, req, sender)
	default:
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusNotFound,
			Body:   []byte("Resource not found"),
		})
	}
}

func (d *Datasource) handleHealth(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	result, _ := d.CheckHealth(ctx, nil)
	body, _ := json.Marshal(result)
	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   body,
	})
}

func (d *Datasource) handleBuckets(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	result, err := d.s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusInternalServerError,
			Body:   []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
		})
	}

	buckets := make([]string, 0, len(result.Buckets))
	for _, b := range result.Buckets {
		buckets = append(buckets, aws.ToString(b.Name))
	}

	body, _ := json.Marshal(map[string][]string{"buckets": buckets})
	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   body,
	})
}

func (d *Datasource) handleFiles(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	var filesReq FilesRequest
	if err := json.Unmarshal(req.Body, &filesReq); err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusBadRequest,
			Body:   []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
		})
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(filesReq.Bucket),
	}
	// Only set prefix if it's not empty and not a wildcard pattern
	if filesReq.Prefix != "" && filesReq.Prefix != "*" && !strings.Contains(filesReq.Prefix, "*") {
		input.Prefix = aws.String(filesReq.Prefix)
	}

	result, err := d.s3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusInternalServerError,
			Body:   []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
		})
	}

	files := make([]map[string]interface{}, 0, len(result.Contents))
	for _, obj := range result.Contents {
		files = append(files, map[string]interface{}{
			"key":          aws.ToString(obj.Key),
			"size":         obj.Size,
			"lastModified": obj.LastModified.Format(time.RFC3339),
		})
	}

	body, _ := json.Marshal(map[string]interface{}{
		"files":     files,
		"truncated": result.IsTruncated != nil && *result.IsTruncated,
	})
	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   body,
	})
}

func (d *Datasource) handleQuery(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	var qm QueryModel
	if err := json.Unmarshal(req.Body, &qm); err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusBadRequest,
			Body:   []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
		})
	}

	// Create a backend query from the request
	queryJSON, _ := json.Marshal(qm)
	backendQuery := backend.DataQuery{
		RefID: qm.RefID,
		JSON:  queryJSON,
	}

	// Execute the query
	response := d.query(ctx, backend.PluginContext{}, backendQuery)

	// Convert frames to JSON response
	if response.Error != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusInternalServerError,
			Body:   []byte(fmt.Sprintf(`{"error": "%s"}`, response.Error.Error())),
		})
	}

	// Build response with frames data
	framesData := make([]map[string]interface{}, 0, len(response.Frames))
	for _, frame := range response.Frames {
		frameData := map[string]interface{}{
			"schema": map[string]interface{}{
				"name":   frame.Name,
				"refId":  frame.RefID,
				"fields": make([]map[string]interface{}, 0, len(frame.Fields)),
			},
			"data": map[string]interface{}{
				"values": make([]interface{}, 0, len(frame.Fields)),
			},
		}

		fields := frameData["schema"].(map[string]interface{})["fields"].([]map[string]interface{})
		values := frameData["data"].(map[string]interface{})["values"].([]interface{})

		for _, field := range frame.Fields {
			// Map Go field types to Grafana field type names
			// Mark time column as "time" type
			fieldType := mapFieldTypeToGrafana(field.Type())
			if field.Name == qm.TimeColumn {
				fieldType = "time"
			}
			fields = append(fields, map[string]interface{}{
				"name": field.Name,
				"type": fieldType,
			})

			// Extract field values - convert nanoseconds to milliseconds for time fields
			fieldValues := make([]interface{}, field.Len())
			for i := 0; i < field.Len(); i++ {
				val := field.At(i)
				// Convert nanoseconds to milliseconds for time column
				if field.Name == qm.TimeColumn {
					if intVal, ok := val.(int64); ok {
						// Convert nanoseconds to milliseconds
						fieldValues[i] = intVal / 1000000
					} else {
						fieldValues[i] = val
					}
				} else {
					fieldValues[i] = val
				}
			}
			values = append(values, fieldValues)
		}

		frameData["schema"].(map[string]interface{})["fields"] = fields
		frameData["data"].(map[string]interface{})["values"] = values
		framesData = append(framesData, frameData)
	}

	body, _ := json.Marshal(map[string]interface{}{"frames": framesData})
	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   body,
	})
}

// mapFieldTypeToGrafana maps Go data field types to Grafana field type names
func mapFieldTypeToGrafana(fieldType data.FieldType) string {
	switch fieldType {
	case data.FieldTypeTime:
		return "time"
	case data.FieldTypeInt8, data.FieldTypeInt16, data.FieldTypeInt32, data.FieldTypeInt64,
		data.FieldTypeUint8, data.FieldTypeUint16, data.FieldTypeUint32, data.FieldTypeUint64,
		data.FieldTypeFloat32, data.FieldTypeFloat64:
		return "number"
	case data.FieldTypeBool:
		return "boolean"
	case data.FieldTypeString:
		return "string"
	default:
		return "other"
	}
}

func (d *Datasource) handleSchema(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	var schemaReq SchemaRequest
	if err := json.Unmarshal(req.Body, &schemaReq); err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusBadRequest,
			Body:   []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
		})
	}

	// Get object
	result, err := d.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(schemaReq.Bucket),
		Key:    aws.String(schemaReq.Key),
	})
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusInternalServerError,
			Body:   []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
		})
	}
	defer result.Body.Close()

	// Read schema
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusInternalServerError,
			Body:   []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
		})
	}

	schema, err := parquetReader.GetSchema(data)
	if err != nil {
		return sender.Send(&backend.CallResourceResponse{
			Status: http.StatusInternalServerError,
			Body:   []byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())),
		})
	}

	body, _ := json.Marshal(schema)
	return sender.Send(&backend.CallResourceResponse{
		Status: http.StatusOK,
		Body:   body,
	})
}
