package s3

/*
==== S3 ==========================================================
|			                                                           |
| S3 utility wrappers for the following interactions: 	         |
| - Deleting an object from S3			  	                         |
|	- Uploading a new object to s3                                 |
|																						                     |
==================================================================
*/

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const DefaultContentType = "application/octet-stream"

const (
	USER_PROFILE = "uploads/users/"
)

func BuildS3Key(dir string, authorID string, filename string) string {
	var builder strings.Builder

	builder.WriteString(dir)
	builder.WriteString(authorID)
	builder.WriteString("/")

	escapedFilename := strings.ReplaceAll(filename, " ", "-")

	builder.WriteString(escapedFilename)

	return builder.String()
}

func DeleteFromS3(key string) error {
	err := hasS3Credentials()
	if err != nil {
		return err
	}

	encodedKey := url.PathEscape(key)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	s3Client := s3.NewFromConfig(cfg)

	bucket := os.Getenv("AWS_BUCKET")

	if _, err = s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &encodedKey,
	}); err != nil {
		return err
	}

	return nil
}

func UploadToS3New(fileHeader *multipart.FileHeader, fileData []byte, key string) (string, error) {
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

func UploadToS3(fileHeader *multipart.FileHeader, fileData []byte, userId string) (string, error) {
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

	escapedFileName := strings.ReplaceAll(fileHeader.Filename, " ", "-")
	key := "featured_images/users/" + userId + "/" + escapedFileName

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
