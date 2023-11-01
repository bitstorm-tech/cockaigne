package persistence

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2/log"
)

var Bucket = os.Getenv("DO_SPACES_BUCKET")
var BaseUrl = os.Getenv("DO_SPACES_BASE_URL")
var S3 *s3.Client

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

	S3 = s3.NewFromConfig(cfg)

	log.Infof("S3 init done: region=%s, endpoint=%s, bucket=%s, baseUrl=%s", region, endpoint, Bucket, BaseUrl)
}
