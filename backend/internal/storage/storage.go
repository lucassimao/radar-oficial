package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type SpacesUploader struct {
	Client *minio.Client
	Bucket string
}

func NewSpacesUploader(bucket string) (*SpacesUploader, error) {
	endpoint := os.Getenv("DO_SPACES_ENDPOINT")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init spaces client: %w", err)
	}

	return &SpacesUploader{
		Client: minioClient,
		Bucket: bucket,
	}, nil
}

// UploadFile uploads content from a reader with known size
func (s *SpacesUploader) UploadFile(ctx context.Context, objectKey string, content io.Reader, size int64, contentType string) error {
	_, err := s.Client.PutObject(ctx, s.Bucket, objectKey, content, size, minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: map[string]string{"x-amz-acl": "public-read"},
	})
	if err != nil {
		return fmt.Errorf("failed to upload to Spaces: %w", err)
	}

	log.Printf("✅ Uploaded to Spaces: %s/%s", s.Bucket, objectKey)
	return nil
}
