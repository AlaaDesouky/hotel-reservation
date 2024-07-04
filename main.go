package main

import (
	"context"
	"hotel-reservation/api"
	"hotel-reservation/db"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var fiberConfig = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func init(){
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	uri := os.Getenv("MONGO_DB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}


	// Initializations
	var (
		// Stores
		userStore = db.NewMongoUserStore(client)
		hotelStore = db.NewMongoHotelStore(client)
		roomStore = db.NewMongoRoomStore(client, hotelStore)
		store = &db.Store{
			User: userStore,
			Hotel: hotelStore,
			Room: roomStore,
		}

		// Handlers
		userHandler = api.NewUserHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		roomHandler = api.NewRoomHandler(store)
		authHandler = api.NewAuthHandler(userStore)

		// Api
		app = fiber.New(fiberConfig)
		authV1 = app.Group("/api/auth")
		apiV1 = app.Group("/api/v1", api.JWTAuthentication(userStore))
	)

	// Auth handlers
	authV1.Post("/login",authHandler.HandelAuthenticate)
	authV1.Post("/signup", authHandler.HandleCreateUser)

	// User handlers
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// Hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetHotelRooms)

	// Room handlers
	apiV1.Get("/room", roomHandler.HandleGetRooms)

	listenAddr := os.Getenv("LISTEN_ADDRESS")
	app.Listen(listenAddr)
}
