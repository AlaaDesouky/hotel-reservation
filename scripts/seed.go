package main

import (
	"context"
	"crypto/rand"
	"hotel-reservation/db"
	"hotel-reservation/db/fixtures"
	"log"
	"math/big"
	mRand "math/rand"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		ctx = context.Background()
		mongoDBUrl = os.Getenv("MONGO_DB_URI")
		mongoDBName = os.Getenv("MONGO_DB_NAME")
		roomSizes = []string{"small", "medium", "large"}
	) 

	client, err := mongo.Connect(ctx, options.Client().ApplyURI((mongoDBUrl)))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User: db.NewMongoUserStore(client),
		Hotel: hotelStore,
		Room: db.NewMongoRoomStore(client, hotelStore),
	}

	var wg sync.WaitGroup

	for range 100 {
		wg.Add(1)
		go func(){
			defer wg.Done()

			fixtures.AddUser(store, generateRandomName(7), generateRandomName(7), mRand.Intn(2) == 0)

			hotel := fixtures.AddHotel(store, generateRandomName(7) + " hotel", generateRandomName(7) + " city", mRand.Intn(5)+1, nil)

			var roomWg sync.WaitGroup

			for range 50 {
				roomWg.Add(1)
				go func(){
					defer roomWg.Done()
					fixtures.AddRoom(store,roomSizes[mRand.Intn(3)], mRand.Intn(2) == 0, mRand.Float64() * 100, hotel.ID)
				}()

				roomWg.Wait()
			}
		}()
	}

	wg.Wait()
}

func generateRandomName(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range result {
		randomIndex, _ := rand.Int(rand.Reader, charsetLength)
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result)
}
