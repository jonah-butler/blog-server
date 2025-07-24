package s3

/*
==== S3 Utilities ================================================
|			                                                           |
| Non-exported s3-related, quality-of-life utility 			         |
| functions for interacting with the AWS Go SDK							  	 |
|																						                     |
==================================================================
*/

import (
	"fmt"
	"os"
	"path/filepath"
)

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

func getS3FileURL(bucket, region, key string) string {
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)
}

func hasS3Credentials() error {
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
