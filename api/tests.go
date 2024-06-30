package api

import (
	"context"
	"hotel-reservation/db"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDB struct {
	client *mongo.Client
	*db.Store
} 

func (tdb *testDB) teardown(t *testing.T) {
	dbName := os.Getenv(db.DB_NAME)
	if err := tdb.client.Database(dbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}

	uri := os.Getenv("MONGO_TEST_DB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	
	return &testDB{
		client: client,
		Store: &db.Store{
			User: db.NewMongoUserStore(client),
		},
	}
}