package main

import (
	"flag"
	"hotel-reservation/api"

	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := flag.String("listenAddr", ":8080", "The listen address for the api server")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/users", api.HandleGetUsers)
	apiV1.Get("/users/:id", api.HandleGetUser)
	app.Listen(*listenAddr)
}
