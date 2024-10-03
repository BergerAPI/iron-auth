package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	// Create a new engine
	engine := html.New("./public", ".html")

	// Creating a fiber app with all views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/login", func(ctx *fiber.Ctx) error {
		clientId := ctx.Query("client_id", "")
		redirectUri := ctx.Query("redirect_uri", "")
		state := ctx.Query("state", "")
		success := ctx.Query("success", "true")

		return ctx.Render("login", fiber.Map{
			"ClientId":    clientId,
			"RedirectUri": redirectUri,
			"State":       state,
			"Success":     success,
		})
	})

	app.Post("/login", func(ctx *fiber.Ctx) error {
		// Parsing the form data
		email := ctx.FormValue("email", "")
		password := ctx.FormValue("password", "")

		// Form data also included information about the oauth2 session
		// if it was provided by the query parameters in GET /login
		clientId := ctx.FormValue("client_id", "")
		redirectUri := ctx.FormValue("redirect_uri", "")
		state := ctx.FormValue("state", "")

		if email != "test@niclas.lol" && password != "test123" {
			return ctx.Redirect("/login?success=false&client_id="+clientId+"&redirect_uri="+redirectUri+"&state="+state, 302)
		}

		return ctx.SendString(email + " " + password)
	})

	err := app.Listen(":3000")
	if err == nil {
		println("Something went wrong with starting up the server.")
	}
}
