package persistence

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gofiber/fiber/v2/log"
)

var bucket = os.Getenv("DO_SPACES_BUCKET")
var dealsFolder = "deals"
var s3Client *s3.Client

func InitS3() {
	key := os.Getenv("DO_SPACES_KEY")
	secret := os.Getenv("DO_SPACES_SECRET")
	region := os.Getenv("DO_SPACES_REGION")
	endpoint := os.Getenv("DO_SPACES_ENDPOINT")

	creds := credentials.NewStaticCredentialsProvider(key, secret, "")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: endpoint,
		}, nil
	})
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Panicf("can't create config for DigitalOcean space: %v", err)
	}

	s3Client = s3.NewFromConfig(cfg)

	baseUrl := os.Getenv("DO_SPACES_BASE_URL")

	log.Infof("S3 init done: region=%s, endpoint=%s, bucket=%s, baseUrl=%s", region, endpoint, bucket, baseUrl)
}

func UploadDealImage(image multipart.FileHeader, dealId string, prefix string) error {
	tokens := strings.Split(image.Filename, ".")
	fileExtension := tokens[len(tokens)-1]
	contentType := image.Header.Get("Content-Type")
	if len(contentType) == 0 {
		contentType = strings.ToLower("image/" + fileExtension)
	}
	key := fmt.Sprintf("%s/%s/%s%d.%s", dealsFolder, dealId, prefix, time.Now().UnixMilli(), fileExtension)
	file, err := image.Open()
	if err != nil {
		return err
	}

	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        file,
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})

	return err
}
