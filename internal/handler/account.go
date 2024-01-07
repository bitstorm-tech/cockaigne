package handler

import (
	"errors"
	"net/http"

	"github.com/bitstorm-tech/cockaigne/internal/model"
	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterAccountHandlers(e *echo.Echo) {
	e.GET("/profile-image/:accountId", getProfileImage)
	e.GET("/settings", openSettings)
	e.GET("/settings-user-common", getUserCommonsSettings)
	e.GET("/settings-profile-image", getProfileImageSettings)
	e.POST("/settings", updateAccount)
	e.POST("/profile-image-update", updateProfileImage)
	e.POST("/api/accounts/filter", updateFilter)
	e.POST("/api/accounts/use-location-service", updateUseLocationService)
}

func updateProfileImage(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	deleteImage := c.FormValue("delete-image") == "true"
	if deleteImage {
		err := service.DeleteProfileImage(userId.String())
		if err != nil {
			return view.RenderAlert("Profilbild kann nicht gelöscht werden, bitte versuche es später noch einmal.", c)
		}
	}

	file, err := c.FormFile("profile-image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		zap.L().Sugar().Error("can't get image from request ", err)
		return view.RenderAlert("Profilbild kann nicht gespeichert werden, bitte versuche es später noch einmal.", c)

	}

	if file != nil && !deleteImage {
		_, err = service.SaveProfileImage(userId.String(), file)
		if err != nil {
			zap.L().Sugar().Error("can't save profile image: ", err)
			return view.RenderAlert("Profilbild kann nicht gespeichert werden, bitte versuche es später noch einmal.", c)
		}
	}

	return view.RenderToast("Profilbild erfolgreich geändert", c)
}

func updateAccount(c echo.Context) error {
	username := c.FormValue("username")
	usernameExists := service.UsernameExists(username)

	if usernameExists {
		return view.RenderAlert("Der Benutzername ist leider schon vergeben.", c)
	}

	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	err = service.UpdateUsername(userId.String(), username)
	if err != nil {
		zap.L().Sugar().Error("can't update username: ", err)
		return view.RenderAlert("Dein Account kann moment nicht geändert werden, bitte versuche es später nochmal.", c)
	}

	return view.RenderToast("Benutzername erfolgreich geändert", c)
}

func getProfileImageSettings(c echo.Context) error {
	user, err := service.ParseUser(c)
	if err != nil {
		return redirect.Login(c)
	}

	imageUrl, err := service.GetProfileImage(user.ID.String())
	if err != nil {
		zap.L().Sugar().Errorf("can't get profile image of user '%s': %v", user.ID, err)
		return view.RenderAlert("Kann Profilbild nicht laden, bitte versuche es später nochmal.", c)
	}

	return view.Render(view.ProfileImageSettings(imageUrl, user.IsDealer), c)
}

func getUserCommonsSettings(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	acc, err := service.GetAccount(userId.String())
	if err != nil {
		zap.L().Sugar().Errorf("can't get user account by id '%s': %v", userId, err)
		return view.RenderAlert("Deine Einstellungen konnten gerade nicht geladen werden, bitte versuche es später nochmal.", c)
	}

	return view.Render(view.CommonUserSettings(acc.Username, acc.Email), c)
}

func openSettings(c echo.Context) error {
	return view.Render(view.Settings(), c)
}

func getProfileImage(c echo.Context) error {
	accountId := c.Param("accountId")
	isDealer := c.QueryParam("dealer") == "true"

	imageUrl, err := service.GetProfileImage(accountId)
	if err != nil {
		zap.L().Sugar().Errorf("can't get profile image for account '%s': %v", accountId, err)
		return view.Render(view.ProfileImage("", isDealer), c)
	}

	return view.Render(view.ProfileImage(imageUrl, isDealer), c)
}

func updateFilter(c echo.Context) error {
	userId, _ := service.ParseUserId(c)

	updateFilterRequest := model.UpdateFilterRequest{}

	err := c.Bind(&updateFilterRequest)
	if err != nil {
		zap.L().Sugar().Error("can't parse filter update request: ", err)
	}

	if err := service.UpdateSearchRadius(userId, updateFilterRequest.SearchRadiusInMeters); err != nil {
		zap.L().Sugar().Error("can't update accounts search_radius_in_meters: ", err)
		return view.RenderAlert("Fehler beim Verarbeiten der Filteränderung", c)
	}

	if err := service.UpdateSelectedCategories(userId, updateFilterRequest.FavoriteCategoryIds); err != nil {
		zap.L().Sugar().Error("can't update selected categories: ", err)
		return view.RenderAlert("Fehler beim Verarbeiten der Filteränderung", c)
	}

	return nil
}

func updateUseLocationService(c echo.Context) error {
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	useLocationService := c.FormValue("use-location-service") == "on"
	err = service.UpdateUseLocationService(userId.String(), useLocationService)
	if err != nil {
		zap.L().Sugar().Error("can't save use location service: ", err)
		return view.RenderAlert("Kann Einstellung leider nicht speichern, bitte später nochmal versuchen.", c)
	}

	if !useLocationService {
		address := c.FormValue("address")
		point, err := service.GetPositionFromAddressFuzzy(address)
		if err != nil {
			zap.L().Sugar().Errorf("can't find position from address (%s): %v", address, err)
			return view.RenderAlert("Ungültige Adresse, bitte geben Sie eine genauere Adresse an.", c)
		}
		err = service.UpdateLocation(userId.String(), point)
		if err != nil {
			zap.L().Sugar().Errorf("can't update location (%s): %v", address, err)
			return view.RenderAlert("Kann Einstellung leider nicht speichern, bitte später nochmal versuchen.", c)
		}
	}

	return nil
}
