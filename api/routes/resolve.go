package routes

import (
	"github.com/Azathoth-X/url-shorten-go-redis/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveUrl(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	val, err := r.Get(database.Ctx, url).Result()

	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "url no in db",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "no db connection",
		})
	}

	rdb := database.CreateClient(1)
	defer rdb.Close()

	_ = rdb.Incr(database.Ctx, "counter")
	return c.Redirect(val, 301)

}
