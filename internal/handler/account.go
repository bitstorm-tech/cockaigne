package handler

import (
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
	e.GET("/settings-user-profile-image", getUserProfileImageSettings)
	e.POST("/api/accounts/filter", updateFilter)
	e.POST("/api/accounts/use-location-service", updateUseLocationService)
}

func getUserProfileImageSettings(c echo.Context) error {
	return view.Render(view.ProfileImageUserSettings(), c)
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
		c.Logger().Errorf("can't parse filter update request: %v", err)
	}

	if err := service.UpdateSearchRadius(userId, updateFilterRequest.SearchRadiusInMeters); err != nil {
		c.Logger().Errorf("can't update accounts search_radius_in_meters: %v", err)
		return view.RenderAlert("Fehler beim Verarbeiten der Filteränderung", c)
	}

	if err := service.UpdateSelectedCategories(userId, updateFilterRequest.FavoriteCategoryIds); err != nil {
		c.Logger().Errorf("can't update selected categories: %s", err)
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
		c.Logger().Errorf("can't save use location service: %v", err)
		return view.RenderAlert("Kann Einstellung leider nicht speichern, bitte später nochmal versuchen.", c)
	}

	if !useLocationService {
		address := c.FormValue("address")
		point, err := service.GetPositionFromAddressFuzzy(address)
		if err != nil {
			c.Logger().Errorf("can't find position from address (%s): %v", address, err)
			return view.RenderAlert("Ungültige Adresse, bitte geben Sie eine genauere Adresse an.", c)
		}
		err = service.UpdateLocation(userId.String(), point)
		if err != nil {
			c.Logger().Errorf("can't update location (%s): %v", address, err)
			return view.RenderAlert("Kann Einstellung leider nicht speichern, bitte später nochmal versuchen.", c)
		}
	}

	return nil
}
