package dealer

import (
	"fmt"
	"github.com/bitstorm-tech/cockaigne/internal/persistence"
	"mime/multipart"
	"strings"
	"time"
)

var folder = "dealer"

func SaveDealerImage(dealerId string, image *multipart.FileHeader) (string, error) {
	tokens := strings.Split(image.Filename, ".")
	fileExtension := tokens[len(tokens)-1]
	path := fmt.Sprintf("%s/%s/%d.%s", folder, dealerId, time.Now().UnixMilli(), fileExtension)
	err := persistence.UploadImage(path, image)
	if err != nil {
		return "", err
	}

	imageUrl := persistence.BaseUrl + "/" + path

	return imageUrl, nil
}

func GetDealerImageUrls(dealerId string) ([]string, error) {
	imageUrls, err := persistence.GetImageUrls(folder + "/" + dealerId)
	if err != nil {
		return []string{}, err
	}

	return imageUrls, nil
}

func DeleteDealerImage(imageUrl string) error {
	path := strings.Replace(imageUrl, persistence.BaseUrl+"/", "", -1)
	return persistence.DeleteImage(path)
}
