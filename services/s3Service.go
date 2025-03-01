package services

import (
	"context"
	"fluxton/errs"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"os"
	"strings"
)

type S3Service interface {
	CreateBucket(bucketName string) (*s3.CreateBucketOutput, error)
	ListBuckets(limit int, continuationToken *string) ([]string, *string, error)
	ShowBucket(bucketName string) (*s3.HeadBucketOutput, error)
	DeleteBucket(bucketName string) error
}

type S3ServiceImpl struct {
	client *s3.Client
}

func NewS3Service() (S3Service, error) {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""))),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	// Create an S3 client with the loaded config
	client := s3.NewFromConfig(cfg)

	return &S3ServiceImpl{
		client: client,
	}, nil
}

func (s *S3ServiceImpl) CreateBucket(bucketName string) (*s3.CreateBucketOutput, error) {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
		LocationConstraint: types.BucketLocationConstraint(os.Getenv("AWS_REGION")),
	}

	createdBucket, err := s.client.CreateBucket(context.Background(), input)
	if err != nil {
		return nil, s.TransformError(fmt.Errorf("createBucket: %q, %v", bucketName, err))
	}

	if !s.BucketExists(bucketName) {
		return nil, errs.NewBadRequestError(fmt.Sprintf("failed to confirm bucket %q exists", bucketName))
	}

	return createdBucket, nil
}

func (s *S3ServiceImpl) BucketExists(bucketName string) bool {
	_, err := s.client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return false
	}

	return true
}

func (s *S3ServiceImpl) ListBuckets(limit int, continuationToken *string) ([]string, *string, error) {
	input := &s3.ListBucketsInput{}

	if limit > 0 {
		input.MaxBuckets = aws.Int32(int32(limit))
	}

	if continuationToken != nil {
		input.ContinuationToken = continuationToken
	}

	resp, err := s.client.ListBuckets(context.Background(), input)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to list buckets: %w", err)
	}

	var bucketNames []string
	for _, bucket := range resp.Buckets {
		if bucket.Name != nil {
			bucketNames = append(bucketNames, *bucket.Name)
		}
	}

	var nextToken *string
	if resp.ContinuationToken != nil {
		nextToken = resp.ContinuationToken
	}

	return bucketNames, nextToken, nil
}

func (s *S3ServiceImpl) ShowBucket(bucketName string) (*s3.HeadBucketOutput, error) {
	resp, err := s.client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to show bucket %q, %v", bucketName, err)
	}

	return resp, nil
}

func (s *S3ServiceImpl) DeleteBucket(bucketName string) error {
	_, err := s.client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("unable to delete bucket %q, %v", bucketName, err)
	}

	return nil
}

func (s *S3ServiceImpl) TransformError(err error) error {
	if err == nil {
		return nil
	}

	errorString := err.Error()

	if strings.Contains(errorString, "BucketAlreadyOwnedByYou") {
		return errs.NewNotFoundError("s3.error.bucketAlreadyOwned")
	}

	if strings.Contains(errorString, "BucketAlreadyExists") {
		return errs.NewBadRequestError("s3.error.bucketAlreadyExists")
	}

	if strings.Contains(errorString, "NoSuchBucket") {
		return errs.NewNotFoundError("s3.error.bucketNotFound")
	}

	return err
}
