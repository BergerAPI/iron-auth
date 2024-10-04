package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/url"
)

func Authorize(ctx *fiber.Ctx) error {
	clientId := ctx.Query("client_id", "")
	redirectUri := ctx.Query("redirect_uri", "")
	responseType := ctx.Query("response_type", "")
	state := ctx.Query("state", "")

	// Checking whether the user is logged in
	userId, ok := ctx.Locals("user").(string)

	if !ok {
		return ctx.Redirect(fmt.Sprintf("/login?client_id=%s&redirect_uri=%s&state=%s", clientId, redirectUri, state))
	}

	// [RFC6749] 4.1.1 client_id REQUIRED; when redirect uri is passed
	if clientId == "" && redirectUri != "" {
		return ctx.Redirect(fmt.Sprintf("%s?error=%s", redirectUri, "invalid_request"))

	}

	// [RFC6749] 4.1.1 client_id REQUIRED; when redirect uri is not passed
	// (No information at all, return an error)
	if clientId == "" && redirectUri == "" {
		return ctx.JSON(fiber.Map{"error": "invalid_request"})
	}

	// [RFC6749] 3.1.2 The redirection endpoint URI MUST be an absolute URI
	if _, err := url.ParseRequestURI(redirectUri); err != nil {
		return ctx.JSON(fiber.Map{"error": "invalid_request"})
	}

	// [RFC6749] 3.1.1 The value MUST be one of "code" for requesting an authorization code
	if responseType != "code" {
		return ctx.Redirect(fmt.Sprintf("%s?error=%s", redirectUri, "unsupported_response_type"))
	}

	return ctx.SendString(userId)
}
