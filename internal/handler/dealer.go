package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/bitstorm-tech/cockaigne/internal/redirect"
	"github.com/bitstorm-tech/cockaigne/internal/service"
	"github.com/bitstorm-tech/cockaigne/internal/view"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RegisterDealerHandlers(e *echo.Echo) {
	e.GET("/dealer/:dealerId", getDealerPage)
	e.GET("/deal-overview", getOverviewPage)
	e.GET("/templates", getTemplatesPage)
	e.GET("/dealer-images/:id", getDealerImages)
	e.GET("/dealer-ratings/:dealerId", getDealerRatings)
	e.GET("/dealer-rating-modal/:dealerId", getRatingModal)
	e.GET("/dealer-image-zoom-modal/:dealerId", getImageZoomDialog)
	e.GET("/dealer-header-favorite-button/:dealerId", getDealerHeaderFavoriteButton)
	e.GET("/dealer-subscription-summary", getDealerSubscriptionSummary)
	e.POST("/dealer-rating/:dealerId", createDealerRating)
	e.POST("/dealer-images", addDealerImage)
	e.POST("/dealer-favorite-toggle/:dealerId", toggleDealerFavorite)
	e.DELETE("/dealer-images", deleteDealerImage)
	e.DELETE("/dealer-rating/:dealerId", deleteDealerRating)
}

func getDealerSubscriptionSummary(c echo.Context) error {
	dealerId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	hasSub, err := service.HasActiveSubscription(dealerId.String())
	if err != nil {
		zap.L().Sugar().Info("can't check if dealer has active subscription: ", err)
		return view.Render(view.SubscriptionSummary("", "", true), c)
	}

	if !hasSub {
		return view.Render(view.SubscriptionSummary("", "", false), c)
	}

	freeDaysLeft, err := service.GetFreeDaysLeftFromSubscription(dealerId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get free days left from subscription: ", err)
		return view.Render(view.SubscriptionSummary("", "", true), c)
	}

	endDateString, err := service.GetSubscriptionPeriodEndDate(dealerId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get subscription period end date: ", err)
		return view.Render(view.SubscriptionSummary("", "", true), c)
	}

	freeDaysLeftString := fmt.Sprintf("%d", freeDaysLeft)

	return view.Render(view.SubscriptionSummary(freeDaysLeftString, endDateString, false), c)
}

func getDealerHeaderFavoriteButton(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealerId := c.Param("dealerId")
	isFavorite, err := service.IsFavorite(dealerId, userId.String())
	if err != nil {
		zap.L().Sugar().Error("can't check if dealer is favorite: ", err)
	}

	return view.Render(view.DealerHeaderFavoriteButton(dealerId, isFavorite), c)
}

func toggleDealerFavorite(c echo.Context) error {
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	if user.IsBasicUser {
		return view.RenderInfo("Diese Funktion steht nur PRO-Mitglieder zur Verfügung.", c)
	}

	dealerId := c.Param("dealerId")

	isFavorite, err := service.ToggleDealerFavorite(dealerId, user.ID.String())
	if err != nil {
		zap.L().Sugar().Error("can't toggle dealer favorite: ", err)
		return view.RenderAlert("Kann favorisierte Dealer momentan nicht speichern, bitte versuche es später nochmal.", c)
	}

	return view.Render(view.DealerHeaderFavoriteButton(dealerId, isFavorite), c)
}

func getImageZoomDialog(c echo.Context) error {
	dealerId := c.Param("dealerId")
	startIndex, err := strconv.Atoi(c.QueryParam("index"))
	if err != nil {
		startIndex = 0
	}

	imageUrls, err := service.GetDealerImageUrls(dealerId)
	if err != nil {
		zap.L().Sugar().Error("can't get dealer images: ", err)
		return view.RenderAlert("Kann Dealer Bilder momentan nicht laden, bitte versuche es später nochmal.", c)
	}

	return view.Render(view.ImageZoomModal(imageUrls, startIndex), c)
}

func getDealerPage(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealerId := c.Param("dealerId")
	acc, err := service.GetAccount(dealerId)

	if err != nil {
		zap.L().Sugar().Errorf("can't find dealer (%s): %+v", dealerId, err)
		return c.NoContent(http.StatusNotFound)
	}

	category, err := service.GetCategory(int(acc.DefaultCategory.Int32))
	if err != nil {
		zap.L().Sugar().Error("can't get default category for dealer (%s): %+v", dealerId, err)
	}

	googleMapsLink := fmt.Sprintf(
		"https://maps.google.com/?q=%s %s, %d %s",
		acc.Street.String,
		acc.HouseNumber.String,
		acc.ZipCode.Int32,
		acc.City.String,
	)

	isOwner := dealerId == userId.String()
	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.Dealer(acc, category.Name, isOwner, googleMapsLink, lang), c)
}

func getOverviewPage(c echo.Context) error {
	dealerId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	freeDaysLeft, err := service.GetFreeDaysLeftFromSubscription(dealerId.String())
	freeDaysLeftString := ""
	if err != nil {
		zap.L().Sugar().Error("can't get free days left from subscription: ", err)
	}
	freeDaysLeftString = fmt.Sprintf("%d", freeDaysLeft)
	if freeDaysLeft < 0 {
		freeDaysLeftString = "0"
	}

	periodEndDate, err := service.GetSubscriptionPeriodEndDate(dealerId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get period end date from subscription: ", err)
		periodEndDate = "01.01.3000"
	}

	return view.Render(view.DealsOverview(dealerId.String(), freeDaysLeftString, periodEndDate), c)
}

func getTemplatesPage(c echo.Context) error {
	return view.Render(view.Templates(), c)
}

func getDealerImages(c echo.Context) error {
	dealerId := c.Param("id")
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	imageUrls, err := service.GetDealerImageUrls(dealerId)
	if err != nil {
		zap.L().Sugar().Error("can't get dealer image urls: ", err)
	}

	isOwner := dealerId == userId.String()
	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.DealerImages(imageUrls, isOwner, dealerId, lang), c)
}

func addDealerImage(c echo.Context) error {
	dealerId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	file, err := c.FormFile("image")
	if err != nil {
		zap.L().Sugar().Error("can't get image from post dealer image request: ", err)
	}

	imageUrl, err := service.SaveDealerImage(dealerId.String(), file)
	if err != nil {
		zap.L().Sugar().Error("can't save dealer image: ", err)
	}

	img := fmt.Sprintf("<img src='%s', alt='Dealer image' class='h-36 w-full object-cover' />", imageUrl)

	return c.String(http.StatusOK, img)
}

func deleteDealerImage(c echo.Context) error {
	dealerId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	imageUrl := c.QueryParam("image-url")

	if !strings.Contains(imageUrl, dealerId.String()) {
		zap.L().Sugar().Errorf("not allowed to delete image -> (dealer=%s, imageUrl=%s)", dealerId.String(), imageUrl)
		return nil
	}

	err = service.DeleteDealerImage(imageUrl)
	if err != nil {
		zap.L().Sugar().Error("can't delete dealer image: ", err)
		return view.RenderAlert("Konnte Bild nicht löschen, bitte später nochmal versuchen.", c)
	}

	imageUrls, err := service.GetDealerImageUrls(dealerId.String())
	if err != nil {
		zap.L().Sugar().Error("can't get dealer images: ", err)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.DealerImages(imageUrls, true, dealerId.String(), lang), c)
}

func getDealerRatings(c echo.Context) error {
	dealerId := c.Param("dealerId")
	user, err := service.GetUserFromCookie(c)
	if err != nil {
		zap.L().Sugar().Error("can't get userId: ", err)
	}

	ratings, err := service.GetDealerRatings(dealerId, user.ID.String())
	if err != nil {
		zap.L().Sugar().Errorf("can't get dealer ratings for dealer %s: %+v", dealerId, err)
	}

	isOwner := dealerId == user.ID.String()
	alreadyRated := service.AlreadyRated(dealerId, user.ID.String())
	stars := 0

	for _, rating := range ratings {
		stars += rating.Stars
		rating.CanEdit = rating.UserId == user.ID
	}

	averageRating := float64(stars) / float64(len(ratings))
	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.DealerRatingList(ratings, dealerId, alreadyRated, isOwner, user.IsBasicUser, averageRating, lang), c)
}

func getRatingModal(c echo.Context) error {
	dealerId := c.Param("dealerId")
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	rating, err := service.GetDealerRating(dealerId, userId.String())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return view.RenderAlert("Kann Bewertung momentan nicht ändern, bitte später noch einmal versuchen.", c)
	}

	if rating.Stars < 1 {
		rating.Stars = 5
	}

	canEdit := err == nil
	rating.DealerId, err = uuid.Parse(dealerId)
	if err != nil {
		zap.L().Sugar().Errorf("can't create uuid from string (%s): %v", dealerId, err)
		return view.RenderAlert("Momentan können keine Bewertungen abgegeben werden, bitte versuche es später noch einmal.", c)
	}

	lang := service.GetLanguageFromCookie(c)

	return view.Render(view.DealerRatingModal(rating, canEdit, lang), c)
}

func createDealerRating(c echo.Context) error {
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	dealerId := c.Param("dealerId")
	text := c.FormValue("rating-text")
	starsText := c.FormValue("stars")
	stars, err := strconv.Atoi(starsText)
	if err != nil {
		zap.L().Sugar().Errorf("can't convert stars-text '%s' into int: %+v", starsText, err)
		return view.RenderAlert("Konnte Bewertung nicht speichern, bitte versuche es später noch mal.", c)
	}

	err = service.SaveDealerRating(userId.String(), dealerId, stars, text)
	if err != nil {
		zap.L().Sugar().Error("can't save dealer rating: ", err)
		return view.RenderAlert("Konnte Bewertung nicht speichern, bitte versuche es später noch mal.", c)
	}

	return getDealerRatings(c)
}

func deleteDealerRating(c echo.Context) error {
	dealerId := c.Param("dealerId")
	userId, err := service.GetUserIdFromCookie(c)
	if err != nil {
		return redirect.Login(c)
	}

	err = service.DeleteDealerRating(dealerId, userId.String())
	if err != nil {
		zap.L().Sugar().Error("can't delete dealer rating: ", err)
		return view.RenderAlert("Konnte Bewertung nicht löschen, bitte später noch einmal versuchen.", c)
	}

	return getDealerRatings(c)
}
