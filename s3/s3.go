package s3

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const DefaultContentType = "application/octet-stream"

func UploadToS3(fileHeader *multipart.FileHeader, fileData []byte) (string, error) {
	err := hasS3Credentials()
	if err != nil {
		return "", err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	file := bytes.NewReader(fileData)

	contentType := getContentType(fileHeader.Filename)

	bucketName := os.Getenv("AWS_BUCKET")
	key := "featured_images/" + fileHeader.Filename

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        &bucketName,
		Key:           &key,
		Body:          file,
		ContentLength: &fileHeader.Size,
		ContentType:   contentType,
		ACL:           types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", err
	}

	url := getS3FileURL(bucketName, os.Getenv("AWS_REGION"), key)

	return url, nil
}
