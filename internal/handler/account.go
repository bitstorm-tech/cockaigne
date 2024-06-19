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
	e.GET("/email-change", openEmailChangePage)
	e.GET("/email-change/:code", changeEmail)
	e.POST("/settings", updateAccount)
	e.POST("/settings-dealer-address", updateDealerAddress)
	e.POST("/profile-image-update", updateProfileImage)
	e.POST("/api/send-password-change-email", sendPasswordChangeEmail)
	e.POST("/activate", activateAccount)
	e.POST("/password-change", changePassword)
	e.POST("/api/accounts/filter", updateFilter)
	e.POST("/api/accounts/filter/select-all", selectAllCategories)
	e.POST("/api/accounts/filter/deselect-all", deselectAllCategories)
	e.POST("/api/accounts/location", updateLocation)
	e.POST("/api/accounts/use-location-service", updateUseLocationService)
	e.POST("/api/send-activation-email", sendActivationEmail)
	e.POST("/api/send-email-change-email", sendEmailChangeEmail)
}

func selectAllCategories(c echo.Context) error {
	categories := service.GetCategories()
	var allCategoryIds []int
	for _, c := range categories {
		allCategoryIds = append(allCategoryIds, c.ID)
	}

	return view.Render(view.CategoryList(categories, allCategoryIds), c)
}

func deselectAllCategories(c echo.Context) error {
	categories := service.GetCategories()

	return view.Render(view.CategoryList(categories, []int{}), c)
}

func updateLocation(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return err
	}

	lon, err := strconv.ParseFloat(c.FormValue("lon"), 64)
	if err != nil {
		zap.L().Sugar().Error("can't parse longitude: ", err)
		return err
	}

	lat, err := strconv.ParseFloat(c.FormValue("lat"), 64)
	if err != nil {
		zap.L().Sugar().Error("can't parse latitude: ", err)
		return err
	}

	point := model.Point{
		Lon: lon,
		Lat: lat,
	}

	if user.IsBasicUser {
		service.GetBasicUserFilter(user.ID.String()).Location = point
		return nil
	}

	err = service.UpdateLocation(user.ID.String(), point)
	if err != nil {
		zap.L().Sugar().Error("can't update location: ", err)
		return err
	}

	return nil
}

func changeEmail(c echo.Context) error {
	code := c.Param("code")

	err := service.ChangeEmail(code)
	if err != nil {
		zap.L().Sugar().Error("can't change email: ", err)
		return view.Render(view.EmailChangeResultPage(true), c)
	}

	return view.Render(view.EmailChangeResultPage(false), c)
}

func openEmailChangePage(c echo.Context) error {
	lang := service.GetLanguageFromCookie(c)
	return view.Render(view.EmailChangePage(lang), c)
}

func sendEmailChangeEmail(c echo.Context) error {
	accountId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	baseUrl := service.GetBaseUrl(c)
	newEmail := c.FormValue("email")
	if len(newEmail) == 0 {
		return view.RenderAlertTranslated("alert.invalid_email", c)
	}

	err = service.PrepareEmailChange(accountId.String(), newEmail, baseUrl)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return view.RenderAlertTranslated("alert.email_already_used", c)
		}

		zap.L().Sugar().Error("can't prepare email change: ", err)
		return view.RenderAlertTranslated("alert.can_t_change_email", c)
	}

	return view.RenderInfoTranslated("info.send_activation_email", c)
}

func openSendPasswordChangeEmailPage(c echo.Context) error {
	lang := service.GetLanguageFromCookie(c)
	return view.Render(view.SendPasswordChangeCodePage(lang), c)
}

func changePassword(c echo.Context) error {
	code := c.FormValue("code")
	password := c.FormValue("password")
	passwordRepeat := c.FormValue("password-repeat")

	if password != passwordRepeat {
		return view.RenderAlertTranslated("alert.password_repeat_not_matching", c)
	}

	err := service.ChangePassword(code, password)
	if err != nil {
		zap.L().Sugar().Error("can't change password: ", err)
		return view.RenderAlertTranslated("alert.can_t_change_password", c)
	}

	return view.RenderToast("Passwort erfolgreich geändert.", c)
}

func openPasswordChangePage(c echo.Context) error {
	code := c.Param("code")
	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.PasswordChangePage(code, lang), c)
}

func sendPasswordChangeEmail(c echo.Context) error {
	accountIdString := ""
	accountId, err := service.GetUserIdFromCookie(c)
	if err == nil {
		accountIdString = accountId.String()
	}
	email := c.FormValue("email")
	if len(accountIdString) == 0 && len(email) == 0 {
		return view.RenderAlertTranslated("alert.provide_email", c)
	}

	baseUrl := service.GetBaseUrl(c)
	err = service.PreparePasswordChange(email, accountIdString, baseUrl)
	if err != nil {
		zap.L().Sugar().Error("can't change password: ", err)
		return view.RenderAlertTranslated("alert.can_t_change_password", c)
	}

	return view.RenderInfoTranslated("info.send_change_pw_email", c)
}

func sendActivationEmail(c echo.Context) error {
	email := c.FormValue("email")
	baseUrl := service.GetBaseUrl(c)

	err := service.SendAccountActivationMail(email, baseUrl)
	if err != nil {
		zap.L().Sugar().Error("can't send account activation email: ", err)
		return view.RenderAlertTranslated("alert.can_t_send_activation_email", c)
	}

	return nil
}

func activateAccount(c echo.Context) error {
	codeString := c.FormValue("code")
	code, err := strconv.Atoi(codeString)
	if err != nil {
		zap.L().Sugar().Error("can't convert activation code '%s' to a number: %+v", codeString, err)
		return view.RenderAlertTranslated("alert.invalid_activation_code", c)
	}

	err = service.ActivateAccount(code)
	if err != nil {
		zap.L().Sugar().Error("can't activate account: ", err)
		return view.RenderAlertTranslated("alert.can_t_activate_account", c)
	}

	return redirect.Login(c)
}

func updateDealerAddress(c echo.Context) error {
	dealerId, err := service.GetUserIdFromCookie(c)
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
		return view.RenderAlertTranslated("alert.can_t_update_address", c)
	}

	err = service.UpdateDealerAddress(dealerId.String(), street, houseNumber, city, int32(zip))
	if err != nil {
		zap.L().Sugar().Error("can't update dealer address: ", err)
		return view.RenderAlertTranslated("alert.can_t_update_address", c)
	}

	return view.RenderToast("Adresse erfolgreich geändert", c)
}

func getDealerAddressSettings(c echo.Context) error {
	dealerId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	acc, err := service.GetAccount(dealerId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get account: ", err)
		return view.RenderAlertTranslated("alert.can_t_load_address", c)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.AddressSettings(acc, lang), c)
}

func getDealerCommonSettings(c echo.Context) error {
	dealerId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	account, err := service.GetAccount(dealerId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get account: ", err)
		return view.RenderAlertTranslated("alert.can_t_load_settings", c)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.CommonDealerSettings(account, lang), c)
}

func updateProfileImage(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	deleteImage := c.FormValue("delete-image") == "true"
	if deleteImage {
		err := service.DeleteProfileImage(userId.String())
		if err != nil {
			return view.RenderAlertTranslated("alert.can_t_delet_profile_image", c)
		}
	}

	file, err := c.FormFile("profile-image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		zap.L().Sugar().Error("can't get image from request ", err)
		return view.RenderAlertTranslated("alert.can_t_save_profile_image", c)

	}

	if file != nil && !deleteImage {
		_, err = service.SaveProfileImage(userId.String(), file)
		if err != nil {
			zap.L().Sugar().Error("can't save profile image: ", err)
			return view.RenderAlertTranslated("alert.can_t_save_profile_image", c)
		}
	}

	return view.RenderToast("Profilbild erfolgreich geändert", c)
}

func updateAccount(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	username := c.FormValue("username")
	usernameExists := service.UsernameExists(userId.String(), username)

	if usernameExists {
		return view.RenderAlertTranslated("alert.username_already_exists", c)
	}

	err = service.UpdateUsername(userId.String(), username)
	if err != nil {
		zap.L().Sugar().Error("can't update username: ", err)
		return view.RenderAlertTranslated("alert.can_t_change_username", c)
	}

	phone := c.FormValue("phone")
	err = service.UpdatePhone(userId.String(), phone)
	if err != nil {
		zap.L().Sugar().Error("can't update phone: ", err)
		return view.RenderAlertTranslated("alert.can_t_change_phone", c)
	}

	taxId := c.FormValue("tax-id")
	err = service.UpdateTaxId(userId.String(), taxId)
	if err != nil {
		zap.L().Sugar().Error("can't update tax ID: ", err)
		return view.RenderAlertTranslated("alert.can_t_chante_taxid", c)
	}

	categoryIdString := c.FormValue("category")
	categoryId, err := strconv.Atoi(categoryIdString)
	if err != nil {
		zap.L().Sugar().Error("can't convert cagetory ID from string to int: ", err)
		return view.RenderAlertTranslated("alert.can_t_change_industry", c)
	}

	err = service.UpdateDefaultCategory(userId.String(), categoryId)
	if err != nil {
		zap.L().Sugar().Error("can't update tax ID: ", err)
		return view.RenderAlertTranslated("alert.can_t_change_industry", c)
	}

	return view.RenderToast("Einstellungen erfolgreich geändert", c)
}

func getProfileImageSettings(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	imageUrl, err := service.GetProfileImage(user.ID.String())
	if err != nil {
		zap.L().Sugar().Errorf("can't get profile image of user '%s': %v", user.ID, err)
		return view.RenderAlertTranslated("alert.can_t_load_profile_image", c)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.ProfileImageSettings(imageUrl, user.IsDealer, lang), c)
}

func getUserCommonsSettings(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	acc, err := service.GetAccount(userId.String())
	if err != nil {
		zap.L().Sugar().Errorf("can't get user account by id '%s': %v", userId, err)
		return view.RenderAlertTranslated("alert.can_t_load_settings", c)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.CommonUserSettings(acc.Username, acc.Email, lang), c)
}

func openSettings(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.Settings(user.IsDealer, lang), c)
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
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	updateFilterRequest := model.UpdateFilterRequest{}

	err = c.Bind(&updateFilterRequest)
	if err != nil {
		zap.L().Sugar().Error("can't parse filter update request: ", err)
		return view.RenderAlertTranslated("alert.can_t_save_filter", c)
	}

	if user.IsBasicUser {
		basicUserFilter := service.GetBasicUserFilter(user.ID.String())
		basicUserFilter.SearchRadiusInMeters = updateFilterRequest.SearchRadiusInMeters
		basicUserFilter.SelectedCategories = updateFilterRequest.FavoriteCategoryIds
		return nil
	}

	if err := service.UpdateSearchRadius(user.ID, updateFilterRequest.SearchRadiusInMeters); err != nil {
		zap.L().Sugar().Error("can't update accounts search_radius_in_meters: ", err)
		return view.RenderAlertTranslated("alert.can_t_save_filter", c)
	}

	if err := service.UpdateSelectedCategories(user.ID, updateFilterRequest.FavoriteCategoryIds); err != nil {
		zap.L().Sugar().Error("can't update selected categories: ", err)
		return view.RenderAlertTranslated("alert.can_t_save_filter", c)
	}

	redirectAfterSave := c.QueryParam("redirect-after-save")

	if len(redirectAfterSave) > 0 {
		c.Response().Header().Set("HX-Redirect", redirectAfterSave)
	}

	return nil
}

func updateUseLocationService(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	useLocationService := c.FormValue("use-location-service") == "on"

	if user.IsBasicUser {
		service.GetBasicUserFilter(user.ID.String()).UseLocationService = useLocationService
		return nil
	}

	err = service.UpdateUseLocationService(user.ID.String(), useLocationService)
	if err != nil {
		zap.L().Sugar().Error("can't save use location service: ", err)
		return view.RenderAlertTranslated("alert.can_t_save_settings", c)
	}

	if !useLocationService {
		address := c.FormValue("address")
		point, err := service.GetPositionFromAddressFuzzy(address)
		if err != nil {
			zap.L().Sugar().Errorf("can't find position from address (%s): %v", address, err)
			return view.RenderAlertTranslated("alert.invalid_address", c)
		}
		err = service.UpdateLocation(user.ID.String(), point)
		if err != nil {
			zap.L().Sugar().Errorf("can't update location (%s): %v", address, err)
			return view.RenderAlertTranslated("alert.can_t_save_settings", c)
		}
	}

	return nil
}
