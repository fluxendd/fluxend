package storage

import (
	"context"
	"fluxend/pkg"
	"fluxend/pkg/errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/guregu/null/v6"
	"github.com/samber/do"
	"io"
	"os"
	"strings"
)

type S3ServiceImpl struct {
	client *s3.Client
}

func NewS3Provider(injector *do.Injector) (Provider, error) {
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

func (s *S3ServiceImpl) CreateContainer(bucketName string) (string, error) {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	}

	input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
		LocationConstraint: types.BucketLocationConstraint(os.Getenv("AWS_REGION")),
	}

	createdBucket, err := s.client.CreateBucket(context.Background(), input)
	if err != nil {
		return "", s.transformError(fmt.Errorf("createBucket: %q, %v", bucketName, err))
	}

	if !s.ContainerExists(bucketName) {
		return "", errors.NewBadRequestError(fmt.Sprintf("failed to confirm bucket %q exists", bucketName))
	}

	return pkg.ConvertPointerToString(createdBucket.Location), nil
}

func (s *S3ServiceImpl) ContainerExists(bucketName string) bool {
	_, err := s.client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	return err == nil
}

func (s *S3ServiceImpl) ListContainers(input ListContainersInput) ([]string, string, error) {
	bucketInput := &s3.ListBucketsInput{}

	if input.Limit > 0 {
		bucketInput.MaxBuckets = aws.Int32(int32(input.Limit))
	}

	if input.Token != "" {
		bucketInput.ContinuationToken = &input.Token
	}

	resp, err := s.client.ListBuckets(context.Background(), bucketInput)
	if err != nil {
		return nil, "", fmt.Errorf("unable to list buckets: %w", err)
	}

	var bucketNames []string
	for _, bucket := range resp.Buckets {
		if bucket.Name != nil {
			bucketNames = append(bucketNames, *bucket.Name)
		}
	}

	var nextToken string
	if resp.ContinuationToken != nil {
		nextToken = pkg.ConvertPointerToString(resp.ContinuationToken)
	}

	return bucketNames, nextToken, nil
}

func (s *S3ServiceImpl) ShowContainer(bucketName string) (*ContainerMetadata, error) {
	resp, err := s.client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to show bucket %q, %v", bucketName, err)
	}

	return &ContainerMetadata{
		Identifier: bucketName,
		Region:     null.StringFrom(pkg.ConvertPointerToString(resp.BucketRegion)),
	}, nil
}

func (s *S3ServiceImpl) DeleteContainer(bucketName string) error {
	_, err := s.client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("unable to delete bucket %q, %v", bucketName, err)
	}

	return nil
}

func (s *S3ServiceImpl) UploadFile(input UploadFileInput) error {
	_, err := s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(input.ContainerName),
		Key:    aws.String(input.FileName),
		Body:   strings.NewReader(string(input.FileBytes)),
	})
	if err != nil {
		return fmt.Errorf("unable to upload file %q, %v", input.FileName, err)
	}

	return nil
}

func (s *S3ServiceImpl) RenameFile(input RenameFileInput) error {
	_, err := s.client.CopyObject(context.Background(), &s3.CopyObjectInput{
		Bucket:     aws.String(input.ContainerName),
		CopySource: aws.String(input.ContainerName + "/" + input.FileName),
		Key:        aws.String(input.NewFileName),
	})
	if err != nil {
		return fmt.Errorf("unable to rename file %q to %q, %v", input.FileName, input.NewFileName, err)
	}

	_, err = s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(input.ContainerName),
		Key:    aws.String(input.FileName),
	})
	if err != nil {
		return fmt.Errorf("unable to delete old file %q, %v", input.FileName, err)
	}

	return nil
}

func (s *S3ServiceImpl) DownloadFile(input FileInput) ([]byte, error) {
	resp, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(input.ContainerName),
		Key:    aws.String(input.FileName),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to download file %q, %v", input.FileName, err)
	}

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %q, %v", input.FileName, err)
	}

	return fileBytes, nil
}

func (s *S3ServiceImpl) DeleteFile(input FileInput) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(input.ContainerName),
		Key:    aws.String(input.FileName),
	})
	if err != nil {
		return fmt.Errorf("unable to delete file %q, %v", input.FileName, err)
	}

	return nil
}

func (s *S3ServiceImpl) transformError(err error) error {
	if err == nil {
		return nil
	}

	errorString := err.Error()

	if strings.Contains(errorString, "BucketAlreadyOwnedByYou") {
		return errors.NewNotFoundError("s3.error.bucketAlreadyOwned")
	}

	if strings.Contains(errorString, "BucketAlreadyExists") {
		return errors.NewBadRequestError("s3.error.bucketAlreadyExists")
	}

	if strings.Contains(errorString, "NoSuchBucket") {
		return errors.NewNotFoundError("s3.error.bucketNotFound")
	}

	return err
}
