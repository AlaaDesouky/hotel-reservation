package api

import (
	"hotel-reservation/db"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-API-Token")
		if token == "" {
			return ErrorUnAuthorized()
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		expires, ok := claims["expires"].(float64)
		if !ok || time.Now().Unix() > int64(expires) {
			return ErrorExpiredToken()
		}

		userId, ok := claims["id"].(string)
		if !ok {
			return ErrorUnAuthorized()
		}

		user, err := userStore.GetUserById(c.Context(), userId)
		if err != nil {
			return ErrorUnAuthorized()
		}

		c.Locals("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrorUnAuthorized()
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrorUnAuthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrorUnAuthorized()
	}

	return claims, nil
}
