package s3

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const DefaultContentType = "application/octet-stream"

func HasS3Credentials() error {
	S3Region := os.Getenv("AWS_REGION")
	S3AccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	S3SecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if S3AccessKey != "" && S3SecretKey != "" && S3Region != "" {
		return nil
	}

	if S3AccessKey == "" || S3SecretKey == "" || S3Region == "" {
		return fmt.Errorf("no s3 credentials founds")
	}

	return nil
}

func UploadToS3(fileHeader *multipart.FileHeader, fileData []byte) (string, error) {
	err := HasS3Credentials()
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

// featuredImageLocation... the url
// featuredImgaeKey... just the file name
// featuredImageTag... idk wtf this is or where it comes from in the old code
func getS3FileURL(bucket, key, region string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
}

func getContentType(filename string) *string {
	ext := filepath.Ext(filename)
	contentType := DefaultContentType

	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	}

	return &contentType
}
