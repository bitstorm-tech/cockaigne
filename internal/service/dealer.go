package service

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"github.com/gofiber/fiber/v2/log"
)

func SaveDealerImage(dealerId string, image *multipart.FileHeader) (string, error) {
	tokens := strings.Split(image.Filename, ".")
	fileExtension := tokens[len(tokens)-1]
	path := fmt.Sprintf("%s/%s/%d.%s", persistence.DealerFolder, dealerId, time.Now().UnixMilli(), fileExtension)
	err := persistence.UploadImage(path, image)
	if err != nil {
		return "", err
	}

	imageUrl := persistence.BaseUrl + "/" + path

	return imageUrl, nil
}

func GetDealerImageUrls(dealerId string) ([]string, error) {
	imageUrls, err := persistence.GetImageUrls(persistence.DealerFolder + "/" + dealerId)
	if err != nil {
		return []string{}, err
	}

	return imageUrls, nil
}

func DeleteDealerImage(imageUrl string) error {
	path := strings.Replace(imageUrl, persistence.BaseUrl+"/", "", -1)
	return persistence.DeleteImage(path)
}

func GetDealerRatings(dealerId string, userId string) ([]model.DealerRating, error) {
	var ratings []model.DealerRating
	err := persistence.DB.Select(
		&ratings,
		"select *, user_id = $1 as canEdit from dealer_ratings_view where dealer_id = $2 order by created",
		userId,
		dealerId,
	)

	return ratings, err
}

func GetDealerRating(dealerId string, userId string) (model.DealerRating, error) {
	var rating model.DealerRating
	err := persistence.DB.Get(
		&rating,
		"select * from dealer_ratings_view where dealer_id = $1 and user_id = $2",
		dealerId,
		userId,
	)

	return rating, err
}

func AlreadyRated(dealerId string, userId string) bool {
	rated := true
	err := persistence.DB.Get(
		&rated,
		"select exists(select user_id from dealer_ratings where dealer_id = $1 and user_id = $2)",
		dealerId,
		userId,
	)

	if err != nil {
		log.Errorf("can't check if dealer %s was already rated by user %s", dealerId, userId)
	}

	return rated
}

func SaveDealerRating(userId string, dealerId string, stars int, ratingText string) error {
	query := `
insert into dealer_ratings (user_id, dealer_id, stars, text) values ($1, $2, $3, $4) 
on conflict(user_id, dealer_id) do update set stars = $3, text = $4
`
	_, err := persistence.DB.Exec(
		query,
		userId,
		dealerId,
		stars,
		ratingText,
	)

	return err
}

func DeleteDealerRating(dealerId string, userId string) error {
	_, err := persistence.DB.Exec(
		"delete from dealer_ratings where dealer_id = $1 and user_id = $2",
		dealerId,
		userId,
	)

	return err
}

func ToggleDealerFavorite(dealerId string, userId string) (bool, error) {
	exists, err := IsFavorite(dealerId, userId)
	if err != nil {
		return false, err
	}

	if exists {
		_, err = persistence.DB.Exec("delete from favorite_dealers where user_id = $1 and dealer_id = $2", userId, dealerId)
	} else {
		_, err = persistence.DB.Exec("insert into favorite_dealers (user_id, dealer_id) values ($1, $2)", userId, dealerId)
	}

	return !exists, err
}

func IsFavorite(dealerId string, userId string) (bool, error) {
	exists := false
	err := persistence.DB.Get(
		&exists,
		"select exists(select user_id from favorite_dealers where user_id = $1 and dealer_id = $2)",
		userId,
		dealerId,
	)

	return exists, err
}
