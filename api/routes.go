package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"pgdb-test/db"
	"pgdb-test/util"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	jwtSecret = os.Getenv("JWT_SECRET")
	ctx       = context.Background()
	client    = db.NewClient()
	api       = fiber.New()
)

func Run() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	err = client.Prisma.Connect()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		err = client.Prisma.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	users := api.Group("/users")
	messages := api.Group("/messages")

	/////////////////////////////////////

	api.Post("/register", registerUser)
	api.Post("/login", loginUser)

	api.Use(util.Auth(jwtSecret))

	api.Get("/whoami", whoami)

	users.Get("/:username?", getUsers)
	// users.Delete("/:username?", deleteUser)

	messages.Get("/:message?", getMessages)
	messages.Post("/", postMessage)
	messages.Delete("/:messageID", deleteMessage)

	api.All("*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	api.Listen(":3000")
}
