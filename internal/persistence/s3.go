package persistence

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"go.uber.org/zap"
)

var Bucket = os.Getenv("DO_SPACES_BUCKET")
var S3BaseUrl = os.Getenv("DO_SPACES_BASE_URL")
var S3 *s3.Client
var DealerImagesFolder = "dealer-images"
var DealImagesFolder = "deal-images"
var ProfileImagesFolder = "profile-images"

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
		zap.L().Sugar().Panicf("can't create config for DigitalOcean space: %v", err)
	}

	S3 = s3.NewFromConfig(cfg)

	zap.L().Sugar().Infof("S3 init done: region=%s, endpoint=%s, bucket=%s, baseUrl=%s", region, endpoint, Bucket, S3BaseUrl)
}

func UploadImageS3(path string, image *multipart.FileHeader) error {
	tokens := strings.Split(image.Filename, ".")
	fileExtension := tokens[len(tokens)-1]
	contentType := image.Header.Get("Content-Type")
	if len(contentType) == 0 {
		contentType = strings.ToLower("image/" + fileExtension)
	}

	file, err := image.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &Bucket,
		Key:         &path,
		Body:        file,
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})

	return err
}

func GetImageUrlsS3(path string) ([]string, error) {
	output, err := S3.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &Bucket,
		Prefix: &path,
	})

	if err != nil {
		return []string{}, err
	}

	var imageUrls []string
	for _, content := range output.Contents {
		imageUrls = append(imageUrls, fmt.Sprintf("%s/%s", S3BaseUrl, *content.Key))
	}

	return imageUrls, nil
}

func GetImageUrlS3(path string) (string, error) {
	imageUrls, err := GetImageUrlsS3(path)
	if err != nil {
		return "", err
	}

	if len(imageUrls) > 1 {
		zap.L().Sugar().Error("got more then one profile image URL (use first one)")
	}

	if len(imageUrls) == 0 {
		return "", nil
	}

	return imageUrls[0], nil
}

func DeleteImageS3(path string) error {
	if len(path) == 0 {
		return nil
	}

	zap.L().Sugar().Debugf("delete image: %s", path)
	_, err := S3.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: &Bucket,
		Key:    &path,
	})

	return err
}
