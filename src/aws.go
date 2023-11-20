package zodo

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client

func PushToS3(path, objectKey string) error {
	if s3Client == nil {
		s3Config, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return err
		}
		s3Client = s3.NewFromConfig(s3Config)
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(Config.Sync.S3.Bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	return err
}
