package zodo

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func s3Client() *s3.Client {

	s3Config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	return s3.NewFromConfig(s3Config)
}

func PushToS3(path, objectKey string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = s3Client().PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(Config.Sync.S3.Bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	return err
}

func PullFromS3(path, objectKey string) error {
	result, err := s3Client().GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(Config.Sync.S3.Bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}
	defer result.Body.Close()

	f, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}
	_, err = f.Write(body)
	return err
}
