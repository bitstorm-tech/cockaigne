package persistence

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/media"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
	"go.uber.org/zap"
)

var imageKitBaseUrl = os.Getenv("IMAGEKIT_ENDPOINT_URL")
var ik *imagekit.ImageKit
var dealerImagesFolder = "dealer-images"
var dealImagesFolder = "deal-images"
var profileImagesFolder = "profile-images"

func InitImagekit() {
	var err error
	ik, err = imagekit.New()
	if err != nil {
		panic("can't initialize ImageKit client")
	}

	zap.L().Sugar().Info("ImageKit successfully initialized")
}

func uploadImage(folder string, name string, image *multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	useUniqueFileName := false
	params := uploader.UploadParam{
		FileName:          name,
		Folder:            folder,
		UseUniqueFileName: &useUniqueFileName,
	}

	_, err = ik.Uploader.Upload(context.TODO(), file, params)
	if err != nil {
		return "", err
	}

	imageUrl := fmt.Sprintf("%s/%s/%s", imageKitBaseUrl, folder, name)
	_, err = ik.Media.PurgeCache(context.TODO(), media.PurgeCacheParam{
		Url: imageUrl,
	})
	if err != nil {
		zap.L().Sugar().Warn("can't purge image cache: ", err)
	}

	return imageUrl, nil
}

func deleteImage(folder string, name string) error {
	response, err := ik.Media.Files(context.TODO(), media.FilesParam{
		SearchQuery: fmt.Sprintf(`name="%s"`, name),
		Path:        folder,
	})

	if err != nil {
		return err
	}

	if len(response.Data) == 0 {
		return fmt.Errorf("can't find file for deletion %s/%s", folder, name)
	}

	if len(response.Data) > 1 {
		return fmt.Errorf("found multiple files %s/%s", folder, name)
	}

	_, err = ik.Media.DeleteFile(context.Background(), response.Data[0].FileId)

	return err
}

func getImagesFromFolder(path string, transformations string) ([]string, error) {

	imageUrls := []string{}

	response, err := ik.Media.Files(context.TODO(), media.FilesParam{
		Path: path,
	})

	if err != nil {
		return imageUrls, err
	}

	if len(transformations) > 0 {
		transformations = fmt.Sprintf("?tr=%s", transformations)
	}

	for _, file := range response.Data {
		imageUrls = append(imageUrls, fmt.Sprintf("%s/%s/%s%s", imageKitBaseUrl, path, file.Name, transformations))
	}

	return imageUrls, nil
}

func GetDealImageUrls(dealId string) ([]string, error) {
	path := fmt.Sprintf("%s/%s", dealImagesFolder, dealId)
	return getImagesFromFolder(path, "")
}

func GetDealImageUrlsWithTransformations(dealId string, transformations string) ([]string, error) {
	path := fmt.Sprintf("%s/%s", dealImagesFolder, dealId)
	return getImagesFromFolder(path, transformations)
}

func GetDealerImageUrls(dealerId string) ([]string, error) {
	path := fmt.Sprintf("%s/%s", dealerImagesFolder, dealerId)
	return getImagesFromFolder(path, "")
}

func GetDealerImageUrlsWithTransformations(dealerId string, transformations string) ([]string, error) {
	path := fmt.Sprintf("%s/%s", dealerImagesFolder, dealerId)
	return getImagesFromFolder(path, transformations)
}

func GetProfileImageUrl(accountId string) (string, error) {
	return GetProfileImageUrlWithTransformations(accountId, "")
}

func GetProfileImageUrlWithTransformations(accountId string, transformations string) (string, error) {
	response, err := ik.Media.Files(context.TODO(), media.FilesParam{
		Path:        profileImagesFolder,
		SearchQuery: fmt.Sprintf(`name:"%s"`, accountId),
	})
	if err != nil {
		return "", err
	}

	if len(response.Data) != 1 {
		return "", nil
	}

	if len(transformations) > 0 {
		transformations = fmt.Sprintf("?tr=%s", transformations)
	}

	return fmt.Sprintf("%s/%s/%s%s", imageKitBaseUrl, profileImagesFolder, accountId, transformations), nil
}

func DeleteProfileImage(name string) error {
	return deleteImage(profileImagesFolder, name)
}

func DeleteDealImage(folder string, name string) error {
	return deleteImage(fmt.Sprintf("%s/%s", dealImagesFolder, folder), name)
}

func DeleteDealerImage(folder string, name string) error {
	return deleteImage(fmt.Sprintf("%s/%s", dealerImagesFolder, folder), name)
}

func UploadDealerImage(folder string, name string, image *multipart.FileHeader) (string, error) {
	return uploadImage(fmt.Sprintf("%s/%s", dealerImagesFolder, folder), name, image)
}

func UploadDealImage(folder string, name string, image *multipart.FileHeader) (string, error) {
	return uploadImage(fmt.Sprintf("%s/%s", dealImagesFolder, folder), name, image)
}

func UploadProfileImage(name string, image *multipart.FileHeader) (string, error) {
	return uploadImage(profileImagesFolder, name, image)
}

func CopyDealImages(fromDealId string, toDealId string) error {
	_, err := ik.Media.CopyFolder(context.TODO(), media.CopyFolderParam{
		SourceFolderPath:    fmt.Sprintf("%s/%s/*", dealImagesFolder, fromDealId),
		DestinationPath:     fmt.Sprintf("%s/%s", dealImagesFolder, toDealId),
		IncludeFileVersions: false,
	})

	return err
}
