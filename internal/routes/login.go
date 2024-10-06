package routes

import (
	"github.com/BergerAPI/iron-auth/internal/database"
	"github.com/BergerAPI/iron-auth/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func constructLoginError(error string, clientId string, redirectUri string, state string) string {
	uri, err := utils.CreateURL("/login", map[string]string{
		"error":        error,
		"client_id":    clientId,
		"redirect_uri": redirectUri,
		"state":        state,
	})

	if err == nil {
		return "/"
	}

	return uri
}

func LoginPage(ctx *fiber.Ctx) error {
	clientId := ctx.Query("client_id", "")
	redirectUri := ctx.Query("redirect_uri", "")
	state := ctx.Query("state", "")
	status := ctx.Query("status", "")

	// Checking if the user is logged in; if they are, redirect them away from the login page
	if _, ok := ctx.Locals("user").(string); ok {
		return ctx.Redirect("/")
	}

	return ctx.Render("login", fiber.Map{
		"ClientId":    clientId,
		"RedirectUri": redirectUri,
		"State":       state,
		"Status":      status,
	})
}

func LoginAction(ctx *fiber.Ctx) error {
	// Parsing the form data
	email := ctx.FormValue("email", "")
	password := ctx.FormValue("password", "")

	// Form data also included information about the oauth2 session
	// if it was provided by the query parameters in GET /login
	clientId := ctx.FormValue("client_id", "")
	redirectUri := ctx.FormValue("redirect_uri", "")
	state := ctx.FormValue("state", "")

	// Checking if the user is logged in; if they are, redirect them away from the login action
	if _, ok := ctx.Locals("user").(string); ok {
		return ctx.Redirect("/")
	}

	var user database.User
	if result := database.Instance.Model(database.User{}).First(&user, "email = ?", email); result.Error != nil || user.Password != password {
		return ctx.Redirect(constructLoginError("br", clientId, redirectUri, state))
	}

	// Token and cookie shall become invalid or be removed in 30 days from now
	expiration := time.Now().Add(30 * 24 * time.Hour)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "auth.iron.sh",
		"aud": "iron.sh",
		"id":  user.Id,
		"exp": expiration.Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return ctx.Redirect(constructLoginError("ise", clientId, redirectUri, state))
	}

	// Create and set the cookie for storing the session
	cookie := new(fiber.Cookie)
	cookie.Name = os.Getenv("AUTH_COOKIE")
	cookie.Value = tokenString
	cookie.Expires = expiration
	ctx.Cookie(cookie)

	if redirectUri != "" && clientId != "" {
		uri, err := utils.CreateURL("/oauth/authorize", map[string]string{
			"client_id":     clientId,
			"redirect_uri":  redirectUri,
			"state":         state,
			"response_type": "code",
		})

		if err == nil {
			return ctx.Redirect(uri)
		}
	}

	return ctx.Redirect("/")
}
