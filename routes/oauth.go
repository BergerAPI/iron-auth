package routes

import (
	"github.com/BergerAPI/iron-auth/database"
	"github.com/BergerAPI/iron-auth/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"net/url"
	"os"
)

func constructError(redirectUri string, error string, state string) string {
	uri, err := utils.CreateURL(redirectUri, map[string]string{
		"error": error,
		"state": state,
	})

	if err != nil {
		return "https://auth.iron.sh/"
	}

	return uri
}

func constructLogin(clientId string, redirectUri string, state string) string {
	uri, err := utils.CreateURL("/login", map[string]string{
		"client_id":    clientId,
		"redirect_uri": redirectUri,
		"state":        state,
	})

	if err != nil {
		return "https://auth.iron.sh/"
	}

	return uri
}

func Authorize(ctx *fiber.Ctx) error {
	clientId := ctx.Query("client_id", "")
	redirectUri := ctx.Query("redirect_uri", "")
	responseType := ctx.Query("response_type", "")
	state := ctx.Query("state", "")

	// Checking whether the user is logged in
	userId, ok := ctx.Locals("user").(string)

	if !ok {
		return ctx.Redirect(constructLogin(clientId, redirectUri, state))
	}

	var user database.User
	if result := database.Instance.Model(database.User{}).First(&user, "id = ?", userId); result.Error != nil {
		ctx.ClearCookie(os.Getenv("AUTH_COOKIE"))
		return ctx.Redirect(constructLogin(clientId, redirectUri, state))
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

	// [RFC6749] 4.1.1 client_id REQUIRED; when redirect uri is passed
	if clientId == "" {
		return ctx.Redirect(constructError(redirectUri, "invalid_request", state))
	}

	// [RFC6749] 3.1.1 The value MUST be one of "code" for requesting an authorization code
	if responseType != "code" {
		return ctx.Redirect(constructError(redirectUri, "unsupported_response_type", state))
	}

	// Requesting further information about the client
	var client database.Client
	if result := database.Instance.Model(database.Client{}).First(&client, "id = ?", clientId); result.Error != nil {
		return ctx.Redirect(constructError(redirectUri, "unauthorized_client", state))
	}

	if client.RedirectUri != redirectUri {
		return ctx.Redirect(constructError(redirectUri, "invalid_request", state))
	}

	// Creating the code used for requesting the access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":       "auth.iron.sh",
		"aud":       "iron.sh",
		"client_id": client.Id,
		"user_id":   user.Id,
	})

	// Sign and get the complete encoded token as a string using the secret
	code, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return ctx.Redirect(constructError(redirectUri, "server_error", state))
	}

	successUrl, err := utils.CreateURL(redirectUri, map[string]string{
		"code":  code,
		"state": state,
	})

	if err != nil {
		return ctx.Redirect(constructError(redirectUri, "server_error", state))
	}

	return ctx.Redirect(successUrl)
}
