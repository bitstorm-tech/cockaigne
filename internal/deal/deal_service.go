package deal

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
)

var folder = "deals"
var bucket = persistence.Bucket
var baseUrl = persistence.BaseUrl

type State string

const (
	Active State = "active"
	Past   State = "past"
	Future State = "future"
)

func ToState(state string) State {
	switch strings.ToLower(state) {
	case "active":
		return Active
	case "past":
		return Past
	case "future":
		return Future
	}

	log.Warnf("invalid deal state (%s) -> use active as default", state)

	return Active
}

func GetCategories() []Category {
	var categories []Category
	err := persistence.DB.Select(&categories, "select * from categories where is_active = true")
	if err != nil {
		log.Errorf("can't get categories: %v", err)
	}

	return categories
}

func GetCategory(id int) (Category, error) {
	var category Category
	err := persistence.DB.Get(&category, "select * from categories where id = $1", id)
	if err != nil {
		return Category{}, fmt.Errorf("can't get category (id=%d): %v", id, err)
	}

	return category, nil
}

func SaveDeal(deal Deal) (uuid.UUID, error) {
	var dealId uuid.UUID
	err := persistence.DB.Get(&dealId,
		"insert into deals (dealer_id, title, description, category_id, duration_in_hours, start, template) values ($1, $2, $3, $4, $5, $6, false) returning id",
		deal.DealerId,
		deal.Title,
		deal.Description,
		deal.CategoryId,
		deal.DurationInHours,
		deal.Start,
	)

	if err != nil {
		return uuid.UUID{}, err
	}

	if deal.IsTemplate {
		_, err = persistence.DB.Exec(
			"insert into deals (dealer_id, title, description, category_id, duration_in_hours, start, template) values ($1, $2, $3, $4, $5, $6, true) returning id",
			deal.DealerId,
			deal.Title,
			deal.Description,
			deal.CategoryId,
			deal.DurationInHours,
			deal.Start,
		)
	}

	if err != nil {
		return uuid.UUID{}, err
	}

	return dealId, nil
}

func GetDeal(id string) (Deal, error) {
	var deal Deal
	err := persistence.DB.Get(&deal, "select * from deals where id = $1", id)
	if err != nil {
		return Deal{}, fmt.Errorf("can't get deal from database: %v", err)
	}

	return deal, nil
}

func GetDealsFromView(state State, dealerId *string) ([]DealView, error) {
	if state != Future && state != Active && state != Past {
		return []DealView{}, fmt.Errorf("unknown deal state: %s", state)
	}

	statement := fmt.Sprintf("select *, st_x(location) || ',' || st_y(location) as location from %s_deals_view", state)

	if dealerId != nil {
		statement += fmt.Sprintf(" where dealer_id = '%s'", *dealerId)
	}

	var deals []DealView
	err := persistence.DB.Select(&deals, statement)

	if err != nil {
		return []DealView{}, fmt.Errorf("can't get active deals: %v", err)
	}

	return deals, nil
}

type DealHeader struct {
	ID       uuid.UUID
	Title    string
	Username string
}

func GetDealHeaders() ([]DealHeader, error) {
	var headers []DealHeader
	err := persistence.DB.Select(&headers, "select id, title, username from active_deals_view")
	if err != nil {
		return []DealHeader{}, err
	}

	return headers, nil
}

type DealDetails struct {
	Description string
}

func GetDealDetails(dealId string) (DealDetails, error) {
	var details DealDetails
	err := persistence.DB.Get(&details, "select description from deals where id = $1", dealId)
	if err != nil {
		return DealDetails{}, err
	}

	return details, nil
}

func GetTemplates(dealerId string) ([]Deal, error) {
	var templates []Deal
	err := persistence.DB.Select(&templates, "select * from deals where template = true and dealer_id = $1", dealerId)
	if err != nil {
		return []Deal{}, fmt.Errorf("can't get templates: %v", err)
	}

	return templates, nil
}

func UploadDealImage(image multipart.FileHeader, dealId string, prefix string) error {
	tokens := strings.Split(image.Filename, ".")
	fileExtension := tokens[len(tokens)-1]
	contentType := image.Header.Get("Content-Type")
	if len(contentType) == 0 {
		contentType = strings.ToLower("image/" + fileExtension)
	}
	key := fmt.Sprintf("%s/%s/%s%d.%s", folder, dealId, prefix, time.Now().UnixMilli(), fileExtension)
	file, err := image.Open()
	if err != nil {
		return err
	}

	_, err = persistence.S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &key,
		Body:        file,
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})

	return err
}

func GetDealImageUrls(dealId string) ([]string, error) {
	prefix := fmt.Sprintf("%s/%s", folder, dealId)
	output, err := persistence.S3.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
		Prefix: &prefix,
	})

	if err != nil {
		return []string{}, err
	}

	var imageUrls []string
	for _, content := range output.Contents {
		imageUrls = append(imageUrls, fmt.Sprintf("%s/%s", baseUrl, *content.Key))
	}

	return imageUrls, nil
}
