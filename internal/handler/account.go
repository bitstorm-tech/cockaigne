package handler

import (
	"errors"
	"net/http"
	"strconv"

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
	e.GET("/settings-dealer-common", getDealerCommonSettings)
	e.GET("/settings-dealer-address", getDealerAddressSettings)
	e.GET("/password-change/:code", openPasswordChangePage)
	e.GET("/send-password-change-email", openSendPasswordChangeEmailPage)
	e.POST("/settings", updateAccount)
	e.POST("/settings-dealer-address", updateDealerAddress)
	e.POST("/profile-image-update", updateProfileImage)
	e.POST("/api/send-password-change-email", sendPasswordChangeEmail)
	e.POST("/activate", activateAccount)
	e.POST("/password-change", changePassword)
	e.POST("/api/accounts/filter", updateFilter)
	e.POST("/api/accounts/use-location-service", updateUseLocationService)
	e.POST("/api/send-activation-email", sendActivationEmail)
}

func openSendPasswordChangeEmailPage(c echo.Context) error {
	return view.Render(view.SendPasswordChangeCodePage(), c)
}

func changePassword(c echo.Context) error {
	code := c.FormValue("code")
	password := c.FormValue("password")
	passwordRepeat := c.FormValue("password-repeat")

	if password != passwordRepeat {
		return view.RenderAlert("Das Passwort und die Wiederholung stimmen nicht überein", c)
	}

	err := service.ChangePassword(code, password)
	if err != nil {
		zap.L().Sugar().Error("can't change password: ", err)
		return view.RenderAlert("Das Passwort kann momentan nicht geändert werden. Bitte versuche es später nochmal.", c)
	}

	return view.RenderToast("Passwort erfolgreich geändert.", c)
}

func openPasswordChangePage(c echo.Context) error {
	code := c.Param("code")

	return view.Render(view.PasswordChangePage(code), c)
}

func sendPasswordChangeEmail(c echo.Context) error {
	accountIdString := ""
	accountId, err := service.ParseUserId(c)
	if err == nil {
		accountIdString = accountId.String()
	}
	email := c.FormValue("email")
	if len(accountIdString) == 0 && len(email) == 0 {
		return view.RenderAlert("Bitte E-Mail Adresse angeben.", c)
	}

	baseUrl := service.BuildDomain(c)
	err = service.PasswordChange(email, accountIdString, baseUrl)
	if err != nil {
		zap.L().Sugar().Error("can't change password: ", err)
		return view.RenderAlert("Dein Passwort kann aktuell nicht geändert werden. Bitte versuche es später nochmal.", c)
	}

	return view.RenderInfo("Wir haben dir eine E-Mail zum ändern deines Passworts geschickt.", c)
}

func sendActivationEmail(c echo.Context) error {
	email := c.FormValue("email")
	domain := service.BuildDomain(c)

	err := service.SendAccountActivationMail(email, domain)
	if err != nil {
		zap.L().Sugar().Error("can't send account activation email: ", err)
		return view.RenderAlert("Momentan können keine Aktivierungs-Emails versendet werden. Bitte versuche es später nochmal.", c)
	}

	return nil
}

func activateAccount(c echo.Context) error {
	codeString := c.FormValue("code")
	code, err := strconv.Atoi(codeString)
	if err != nil {
		zap.L().Sugar().Error("can't convert activation code '%s' to a number: %+v", codeString, err)
		return view.RenderAlert("Der angegebene Aktivierungscode ist ungütlig.", c)
	}

	err = service.ActivateAccount(code)
	if err != nil {
		zap.L().Sugar().Error("can't activate account: ", err)
		return view.RenderAlert("Aktivierung des Accounts momentan nicht möglich. Bitte versuche es später nochmal.", c)
	}

	return redirect.Login(c)
}

func updateDealerAddress(c echo.Context) error {
	dealerId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	street := c.FormValue("street")
	houseNumber := c.FormValue("housenumber")
	city := c.FormValue("city")
	zipString := c.FormValue("zip")
	zip, err := strconv.Atoi(zipString)
	if err != nil {
		zap.L().Sugar().Error("can't convert zip string into int: ", err)
		return view.RenderAlert("Kann Adresse momentan nicht aktuallisieren, bitte versuche es später noch einmal.", c)
	}

	err = service.UpdateDealerAddress(dealerId.String(), street, houseNumber, city, int32(zip))
	if err != nil {
		zap.L().Sugar().Error("can't update dealer address: ", err)
		return view.RenderAlert("Kann Adresse momentan nicht aktuallisieren, bitte versuche es später noch einmal.", c)
	}

	return view.RenderToast("Adresse erfolgreich geändert", c)
}

func getDealerAddressSettings(c echo.Context) error {
	dealerId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	acc, err := service.GetAccount(dealerId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get account: ", err)
		return view.RenderAlert("Kann aktuelle Adresse momentan nicht laden, bitte versuche es später nochmal.", c)
	}

	return view.Render(view.AddressSettings(acc), c)
}

func getDealerCommonSettings(c echo.Context) error {
	dealerId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	account, err := service.GetAccount(dealerId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get account: ", err)
		return view.RenderAlert("Kann Einstellungen gerade nicht laden, bitte versuche es später noch einmal.", c)
	}

	return view.Render(view.CommonDealerSettings(account), c)
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
	userId, err := service.ParseUserId(c)
	if err != nil {
		return redirect.Login(c)
	}

	username := c.FormValue("username")
	usernameExists := service.UsernameExists(userId.String(), username)

	if usernameExists {
		return view.RenderAlert("Der Benutzername ist leider schon vergeben.", c)
	}

	err = service.UpdateUsername(userId.String(), username)
	if err != nil {
		zap.L().Sugar().Error("can't update username: ", err)
		return view.RenderAlert("Dein Benutzername kann moment nicht geändert werden, bitte versuche es später nochmal.", c)
	}

	phone := c.FormValue("phone")
	err = service.UpdatePhone(userId.String(), phone)
	if err != nil {
		zap.L().Sugar().Error("can't update phone: ", err)
		return view.RenderAlert("Deine Telefonnummer kann moment nicht geändert werden, bitte versuche es später nochmal.", c)
	}

	taxId := c.FormValue("tax-id")
	err = service.UpdateTaxId(userId.String(), taxId)
	if err != nil {
		zap.L().Sugar().Error("can't update tax ID: ", err)
		return view.RenderAlert("Deine Steuernummer ID kann moment nicht geändert werden, bitte versuche es später nochmal.", c)
	}

	categoryIdString := c.FormValue("category")
	categoryId, err := strconv.Atoi(categoryIdString)
	if err != nil {
		zap.L().Sugar().Error("can't convert cagetory ID from string to int: ", err)
		return view.RenderAlert("Deine Branche kann moment nicht geändert werden, bitte versuche es später nochmal.", c)
	}

	err = service.UpdateDefaultCategory(userId.String(), categoryId)
	if err != nil {
		zap.L().Sugar().Error("can't update tax ID: ", err)
		return view.RenderAlert("Deine Branche kann moment nicht geändert werden, bitte versuche es später nochmal.", c)
	}

	return view.RenderToast("Einstellungen erfolgreich geändert", c)
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
	user, err := service.ParseUser(c)
	if err != nil {
		return redirect.Login(c)
	}

	return view.Render(view.Settings(user.IsDealer), c)
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
