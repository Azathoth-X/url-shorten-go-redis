package routes

import (
	"os"
	"time"

	"github.com/Azathoth-X/url-shorten-go-redis/database"
	"github.com/Azathoth-X/url-shorten-go-redis/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset int           `json:"rate_limit_reset"`
}

func ShortenUrl(c *fiber.Ctx) error {

	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"error": "contparsejson"})
	}

	//rate-limiting

	//check iif url actual url
	if !govalidator.IsUrl(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bruh why no url "})
	}

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": "pls no hack"})
	}

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]

	} else {
		id = body.CustomShort

	}

	r0 := database.CreateClient(0)
	defer r0.Close()

	val, _ := r0.Get(database.Ctx, id).Result()

	if val != "" {
		c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "short in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err := r0.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"err": "internal server go brr"})
	}

	res := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          0,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	res.CustomShort = os.Getenv("DOMAIN") + "/" + id
	return c.Status(fiber.StatusOK).JSON(res)

}
