package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

// AttemptAuthentication This will attempt to authorize
func AttemptAuthentication(ctx *fiber.Ctx) error {
	cookieToken := ctx.Cookies(os.Getenv("AUTH_COOKIE"), "")

	if err := authenticateToken(ctx, cookieToken); err != nil {
		// Removing the token
		if len(cookieToken) != 0 {
			ctx.ClearCookie(os.Getenv("AUTH_COOKIE"))
		}

		return ctx.Next()
	}

	return ctx.Next()
}

func authenticateToken(ctx *fiber.Ctx, tokenString string) error {
	// Parsing the token with JWT
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing the token")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || parsedToken == nil {
		return err
	}

	// Gathering the claims from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok || int64(claims["exp"].(float64)) < time.Now().Local().Unix() {
		return err
	}

	ctx.Locals("user", claims["id"])

	return nil
}
