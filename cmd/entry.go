package main

import (
	"github.com/BergerAPI/iron-auth"
	"github.com/BergerAPI/iron-auth/database"
	"github.com/BergerAPI/iron-auth/routes"
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

	app.Get("/login", middleware.AttemptAuthentication, routes.LoginPage)
	app.Post("/login", middleware.AttemptAuthentication, routes.LoginAction)

	err := app.Listen(":3000")
	if err == nil {
		println("Something went wrong with starting up the server.")
	}
}
