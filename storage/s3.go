package storage

import (
	"context"
	"fmt"
	"time"

	"bytes"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/schollz/progressbar/v3"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

type FileInfo struct {
	Name         string
	Size         int64
	LastModified time.Time
}

func NewS3Client(bucket, endpoint, accessKey, secretKey string) (*S3Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})

	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(creds),
		config.WithRegion("us-east-1"), // default region, can be made configurable
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Client{
		client: client,
		bucket: bucket,
	}, nil
}

func (s *S3Client) ListFiles(prefix string) ([]FileInfo, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: &s.bucket,
		Prefix: &prefix,
	}

	result, err := s.client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	var files []FileInfo
	for _, obj := range result.Contents {
		files = append(files, FileInfo{
			Name:         *obj.Key,
			Size:         *obj.Size,
			LastModified: *obj.LastModified,
		})
	}

	return files, nil
}

func (s *S3Client) UploadFile(localPath, remotePath string, data []byte, progress io.Writer) error {
	size := int64(len(data))
	bar := progressbar.DefaultBytes(size, "uploading")

	// Create a Buffer to hold the data
	buf := bytes.NewReader(data)

	input := &s3.PutObjectInput{
		Bucket:        &s.bucket,
		Key:           &remotePath,
		Body:          buf,
		ContentLength: &size,
	}

	_, err := s.client.PutObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to upload: %v", err)
	}

	// Update progress bar after successful upload
	err = bar.Add64(size)
	if err != nil {
		return err
	}

	return nil
}

func (s *S3Client) DownloadFile(remotePath string, progress io.Writer) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &remotePath,
	}

	result, err := s.client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("failed to close body: %v", err)
		}
	}(result.Body)

	bar := progressbar.DefaultBytes(
		*result.ContentLength,
		"downloading",
	)

	return io.ReadAll(io.TeeReader(result.Body, bar))
}

func (s *S3Client) DeleteFile(remotePath string) error {
	input := &s3.DeleteObjectInput{
		Bucket: &s.bucket,
		Key:    &remotePath,
	}

	_, err := s.client.DeleteObject(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	return nil
}
