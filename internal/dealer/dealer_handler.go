package dealer

import (
	"fmt"
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/auth/jwt"
	"github.com/bitstorm-tech/cockaigne/internal/deal"
	"github.com/bitstorm-tech/cockaigne/internal/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"strconv"
	"strings"
)

func Register(app *fiber.App) {
	app.Get("/dealer/:dealerId", getDealerPage)
	app.Get("/deal-overview", getOverviewPage)
	app.Get("/templates", getTemplatesPage)
	app.Get("/dealer-images/:id", getDealerImages)
	app.Get("/dealer-ratings/:dealerId", getDealerRatings)
	app.Get("/dealer-rating-modal/:dealerId", openRatingModal)
	app.Post("/dealer-rating/:dealerId", createDealerRating)
	app.Post("/dealer-images", addDealerImage)
	app.Delete("/dealer-images", deleteDealerImage)
	app.Delete("/dealer-rating/:dealerId", deleteDealerRating)
}

func getDealerPage(c *fiber.Ctx) error {
	dealerId := c.Params("dealerId")
	acc, err := account.GetAccount(dealerId)

	if err != nil {
		log.Errorf("can't find dealer (%s): %v", dealerId, err)
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}

	category, err := deal.GetCategory(int(acc.DefaultCategory.Int32))

	if err != nil {
		log.Error(err)
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}

	return c.Render("pages/dealer", fiber.Map{"account": acc, "category": category.Name}, "layouts/main")
}

func getOverviewPage(c *fiber.Ctx) error {
	return c.Render("pages/deal-overview", nil, "layouts/main")
}

func getTemplatesPage(c *fiber.Ctx) error {
	userId, err := jwt.ParseUserId(c)
	if err != nil {
		return c.Redirect("/login")
	}

	templates, err := deal.GetTemplates(userId.String())
	if err != nil {
		log.Errorf("can't get templates for dealer: %s", userId.String())
	}

	return c.Render("pages/templates", fiber.Map{"templates": templates}, "layouts/main")
}

func getDealerImages(c *fiber.Ctx) error {
	dealerId := c.Params("id")
	userId, err := jwt.ParseUserId(c)
	if err != nil {
		return c.Redirect("/login")
	}

	imageUrls, err := GetDealerImageUrls(dealerId)
	if err != nil {
		log.Errorf("can't get dealer image urls: %v", err)
	}

	isOwner := len(imageUrls) > 0 && strings.Contains(imageUrls[0], userId.String())

	return c.Render("fragments/dealer/images", fiber.Map{"imageUrls": imageUrls, "isOwner": isOwner})
}

func addDealerImage(c *fiber.Ctx) error {
	dealerId, err := jwt.ParseUserId(c)
	if err != nil {
		return c.Redirect("/login")
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Errorf("can't get image from post dealer image request: %v", err)
	}

	imageUrl, err := SaveDealerImage(dealerId.String(), file)
	if err != nil {
		log.Errorf("can't save dealer image: %v", err)
	}

	img := fmt.Sprintf("<img src='%s', alt='Dealer image' class='h-36 w-full object-cover' />", imageUrl)

	return c.SendString(img)
}

func deleteDealerImage(c *fiber.Ctx) error {
	dealerId, err := jwt.ParseUserId(c)
	if err != nil {
		return c.Redirect("/login")
	}

	imageUrl := c.Query("image-url")

	if !strings.Contains(imageUrl, dealerId.String()) {
		log.Errorf("not allowed to delete image -> (dealer=%s, imageUrl=%s)", dealerId.String(), imageUrl)
		return nil
	}

	err = DeleteDealerImage(imageUrl)
	if err != nil {
		log.Errorf("can't delete dealer image: %v", err)
		return ui.ShowAlert(c, "Konnte Bild nicht löschen, bitte später nochmal versuchen.")
	}

	imageUrls, err := GetDealerImageUrls(dealerId.String())
	if err != nil {
		log.Errorf("can't get dealer images: %v", err)
	}

	return c.Render("fragments/dealer/images", fiber.Map{"imageUrls": imageUrls, "isOwner": true})
}

func getDealerRatings(c *fiber.Ctx) error {
	dealerId := c.Params("dealerId")
	userId, err := jwt.ParseUserId(c)
	if err != nil {
		log.Errorf("can't get userId: %v", err)
	}

	ratings, err := GetRatings(dealerId, userId.String())
	if err != nil {
		log.Errorf("can't get dealer ratings for dealer %s: %v", dealerId, err)
	}

	isOwner := dealerId == userId.String()
	alreadyRated := AlreadyRated(dealerId, userId.String())
	stars := 0

	for _, rating := range ratings {
		stars += rating.Stars
		rating.CanEdit = rating.UserId == userId
	}

	averageRating := float64(stars) / float64(len(ratings))

	return c.Render(
		"fragments/dealer/rating-list",
		fiber.Map{
			"dealerId":      dealerId,
			"ratings":       ratings,
			"isOwner":       isOwner,
			"alreadyRated":  alreadyRated,
			"averageRating": averageRating,
		},
	)
}

func openRatingModal(c *fiber.Ctx) error {
	dealerId := c.Params("dealerId")
	editRating := c.Query("edit") == "true"

	if !editRating {
		rating := Rating{
			Stars: 3,
		}
		return c.Render("fragments/dealer/rating-modal", fiber.Map{"DealerId": dealerId, "Rating": rating})
	}

	userId, err := jwt.ParseUserId(c)
	if err != nil {
		return c.Redirect("/login")
	}

	rating, err := GetRating(dealerId, userId.String())
	if err != nil {
		return ui.ShowAlert(c, "Kann Bewertung momentan nicht ändern, bitte später noch einmal versuchen.")
	}

	return c.Render("fragments/dealer/rating-modal", fiber.Map{"DealerId": dealerId, "Rating": rating, "Edit": editRating})
}

func createDealerRating(c *fiber.Ctx) error {
	userId, err := jwt.ParseUserId(c)
	if err != nil {
		return c.Redirect("/login")
	}

	dealerId := c.Params("dealerId")
	text := c.FormValue("rating-text")
	starsText := c.FormValue("stars")
	stars, err := strconv.Atoi(starsText)
	if err != nil {
		log.Errorf("can't convert stars-text '%s' into int: %v", starsText, err)
		return ui.ShowAlert(c, "Konnte Bewertung nicht speichern, bitte versuche es später noch mal.")
	}

	err = SaveRating(userId.String(), dealerId, stars, text)
	if err != nil {
		log.Errorf("can't save dealer rating: %v", err)
		return ui.ShowAlert(c, "Konnte Bewertung nicht speichern, bitte versuche es später noch mal.")
	}

	return getDealerRatings(c)
}

func deleteDealerRating(c *fiber.Ctx) error {
	dealerId := c.Params("dealerId")
	userId, err := jwt.ParseUserId(c)
	if err != nil {
		return c.Redirect("/login")
	}

	err = DeleteRating(dealerId, userId.String())
	if err != nil {
		log.Errorf("can't delete dealer rating: %v", err)
		return ui.ShowAlert(c, "Konnte Bewertung nicht löschen, bitte später noch einmal versuchen.")
	}

	return getDealerRatings(c)
}
