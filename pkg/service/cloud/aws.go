package cloud

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
)

type awsService struct {
	service    *s3.S3
	bucketName string
}

const (
	filePreSignExpireDuration = time.Hour * 12
)

func NewAWSCloudService(cfg config.Config) (CloudService, error) {

	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.AwsRegion),
		Credentials: credentials.NewStaticCredentials(cfg.AwsAccessKeyID, cfg.AwsSecretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session for aws service : %w", err)
	}

	service := s3.New(session)

	return &awsService{
		service:    service,
		bucketName: cfg.AwsBucketName,
	}, nil
}

func (c *awsService) SaveFile(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {

	file, err := fileHeader.Open()
	if err != nil {
		return "", utils.PrependMessageToError(err, "failed to open file")
	}

	uploadID := uuid.New().String()

	_, err = c.service.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(uploadID),
	})
	if err != nil {
		return "", utils.PrependMessageToError(err, "failed to upload file")
	}

	return uploadID, nil
}
func (c *awsService) GetFileUrl(ctx context.Context, uploadID string) (string, error) {

	req, _ := c.service.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(uploadID),
	})

	url, err := req.Presign(filePreSignExpireDuration)
	if err != nil {
		return "", utils.PrependMessageToError(err, "failed to pre sign url fo uploaded file")
	}

	return url, nil
}
