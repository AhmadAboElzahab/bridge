package storage

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// S3Driver uploads files to any S3-compatible bucket (AWS S3, Cloudflare R2, etc.)
// and returns the file's public URL.
type S3Driver struct {
	client    *s3.Client
	bucket    string
	publicURL string // base public URL without trailing slash
}

// NewS3Driver creates an S3Driver.
//
// For Cloudflare R2:   endpoint = "https://<account-id>.r2.cloudflarestorage.com", region = "auto"
// For AWS S3:          endpoint = "" (uses default), region = "us-east-1" etc.
func NewS3Driver(accessKey, secretKey, region, endpoint, bucket, publicURL string) (*S3Driver, error) {
	if region == "" {
		region = "auto"
	}

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build storage config: %w", err)
	}

	var clientOpts []func(*s3.Options)
	if endpoint != "" {
		clientOpts = append(clientOpts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
			// Path-style is required for R2 and recommended for custom endpoints.
			o.UsePathStyle = true
		})
	}

	return &S3Driver{
		client:    s3.NewFromConfig(cfg, clientOpts...),
		bucket:    bucket,
		publicURL: strings.TrimRight(publicURL, "/"),
	}, nil
}

func (d *S3Driver) Upload(ctx context.Context, folder, ext, contentType string, data []byte) (string, error) {
	key := folder + "/" + uuid.New().String() + ext

	_, err := d.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(d.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return d.publicURL + "/" + key, nil
}
