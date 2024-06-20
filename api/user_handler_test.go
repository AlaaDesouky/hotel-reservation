package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-reservation/types"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func createUserParams() *types.CreateUserParams {
	return &types.CreateUserParams{
		Email: "johndoe@test.com",
		FirstName: "John",
		LastName: "Doe",
		Password: "akldsjasdkljadkfl",
	}
}
func TestGetUsers(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.User)
	app.Post("/", UserHandler.HandlePostUser)
	app.Get("/", UserHandler.HandleGetUsers)

	params := createUserParams()

	b, _ := json.Marshal(params)
	postReq := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	postReq.Header.Add("Content-Type", "application/json")
	_, err := app.Test(postReq)
	if err != nil {
		t.Error(err)
	}

	getReq := httptest.NewRequest("GET", "/", nil)
	getReq.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(getReq)
	if err != nil {
		t.Error(err)
	}

	var users []types.User
	json.NewDecoder(resp.Body).Decode(&users)

	if len(users) != 1 {
		t.Errorf("expected length %d but got %d", 1, len(users))
	}
}

func TestPostUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.User)
	app.Post("/", UserHandler.HandlePostUser)

	params := createUserParams()

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)
	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}
	if len(user.Password) > 0 {
		t.Errorf("expecting the password not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastName %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}

func TestGetUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.User)
	app.Post("/", UserHandler.HandlePostUser)
	app.Get("/:id", UserHandler.HandleGetUser)

	params := createUserParams()

	b, _ := json.Marshal(params)
	postReq := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	postReq.Header.Add("Content-Type", "application/json")
	postResp, err := app.Test(postReq)
	if err != nil {
		t.Error(err)
	}

	var createdUser types.User
	json.NewDecoder(postResp.Body).Decode(&createdUser)
	if len(createdUser.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}

	getReq := httptest.NewRequest("GET", fmt.Sprintf("/%s", createdUser.ID.Hex()), nil)
	getReq.Header.Add("Content-Type", "application/json")

	getResp, err := app.Test(getReq)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(getResp.Body).Decode(&user)

	if len(user.Password) > 0 {
		t.Errorf("expecting the password not to be included in the json response")
	}
	if user.FirstName != params.FirstName {
		t.Errorf("expected firstName %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected lastName %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}