package dealer

import (
	"fmt"
	"github.com/bitstorm-tech/cockaigne/internal/account"
	"github.com/bitstorm-tech/cockaigne/internal/auth/jwt"
	"github.com/bitstorm-tech/cockaigne/internal/deal"
	"github.com/bitstorm-tech/cockaigne/internal/ui"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"strings"
)

func Register(app *fiber.App) {
	app.Get("/dealer/:dealerId", func(c *fiber.Ctx) error {
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
	})

	app.Get("/deal-overview", func(c *fiber.Ctx) error {
		return c.Render("pages/deal-overview", nil, "layouts/main")
	})

	app.Get("/templates", func(c *fiber.Ctx) error {
		userId, err := jwt.ParseUserId(c)
		if err != nil {
			return c.Redirect("/login")
		}

		templates, err := deal.GetTemplates(userId.String())
		if err != nil {
			log.Errorf("can't get templates for dealer: %s", userId.String())
		}

		return c.Render("pages/templates", fiber.Map{"templates": templates}, "layouts/main")
	})

	app.Get("/dealer-images/:id", func(c *fiber.Ctx) error {
		dealerId := c.Params("id")
		imageUrls, err := GetDealerImageUrls(dealerId)
		if err != nil {
			log.Errorf("can't get dealer image urls: %v", err)
		}

		return c.Render("fragments/dealer/images", fiber.Map{"imageUrls": imageUrls})
	})

	app.Post("/dealer-images", func(c *fiber.Ctx) error {
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

		img := fmt.Sprintf("<img src='%s', alt='Dealer image' />", imageUrl)

		return c.SendString(img)
	})

	app.Delete("/dealer-images", func(c *fiber.Ctx) error {
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

		return c.Render("fragments/dealer/images", fiber.Map{"imageUrls": imageUrls})
	})
}
