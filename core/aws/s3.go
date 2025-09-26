package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	S3Client *s3.Client
}

func NewS3(cfg aws.Config) *S3Config {
	return &S3Config{
		S3Client: s3.NewFromConfig(cfg),
	}
}

func (s *S3Config) ListBuckets() ([]string, error) {
	resp, err := s.S3Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	bucketNames := make([]string, len(resp.Buckets))
	for i := 0; i < len(resp.Buckets); i++ {
		bucketNames[i] = *resp.Buckets[i].Name
	}
	return bucketNames, nil
}
