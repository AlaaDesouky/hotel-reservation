package api

import (
	"errors"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(GenericResponse{
		Type: "error",
		Msg: "invalid credentials",
	})
}

func (h *AuthHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if  err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	newUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}

	resp := AuthResponse{
		User: user,
		Token: CreateTokenFromUser(newUser),
	}

	return c.JSON(resp)
}

func (h *AuthHandler) HandelAuthenticate(c *fiber.Ctx) error {
	var params types.AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	fmt.Printf("%+v", user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsValidPassword(user.Password, params.Password) {
		return invalidCredentials(c)
	}

	resp := AuthResponse{
		User: user,
		Token: CreateTokenFromUser(user),
	}

	return c.JSON(resp)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"id": user.ID,
		"email": user.Email,
		"expires": expires,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	
	return tokenStr
}