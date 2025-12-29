package s3

import (
	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// NewS3Client creates a new S3 client with the given configuration
func NewS3Client(endpoint, region, accessKeyID, secretAccessKey string, pathStyle, insecureSkipVerify bool) *s3.Client {
	// Default region
	if region == "" {
		region = "us-east-1"
	}

	// Create credentials
	creds := credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")

	// Create custom HTTP client if needed
	var httpClient *http.Client
	if insecureSkipVerify {
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}

	// Build options
	opts := []func(*s3.Options){
		func(o *s3.Options) {
			o.Region = region
			o.Credentials = creds
			o.UsePathStyle = pathStyle
		},
	}

	if endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	if httpClient != nil {
		opts = append(opts, func(o *s3.Options) {
			o.HTTPClient = httpClient
		})
	}

	return s3.New(s3.Options{}, opts...)
}
