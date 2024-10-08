package main

import (
	"github.com/BergerAPI/iron-auth/internal/database"
	"github.com/BergerAPI/iron-auth/internal/routes"
	"github.com/BergerAPI/iron-auth/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Create a new engine
	engine := html.New("./public", ".html")

	// Creating a fiber app with all views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// Construct a connection to a sqlite database (temporarily a file)
	database.Init("file:db.sqlite")

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello World!")
	})

	app.Get("/login", utils.AttemptAuthentication, routes.LoginPage)
	app.Post("/login", utils.AttemptAuthentication, routes.LoginAction)
	app.Get("/oauth/authorize", utils.AttemptAuthentication, routes.Authorize)
	app.Post("/oauth/token", routes.Token)

	err := app.Listen(":3000")
	if err == nil {
		println("Something went wrong with starting up the server.")
	}
}
