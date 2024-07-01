package db

import (
	"context"
	"fmt"
	"hotel-reservation/types"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const hotelCollection = "hotels"

type HotelStore interface {
	Dropper

	GetHotels(context.Context, Map, *Pagination) ([]*types.Hotel, error)
	GetHotelById(context.Context, string)(*types.Hotel, error)
	CreateHotel(context.Context, *types.Hotel)(*types.Hotel, error)
	UpdateHotel(context.Context, Map, types.UpdateHotelParams) error
	DeleteHotel(context.Context, string) error
}

type MongoHotelStore struct {
	client *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	db_name := os.Getenv(DB_NAME)
	return &MongoHotelStore{
		client: client,
		collection: client.Database(db_name).Collection(hotelCollection),
	}
}

func (s *MongoHotelStore) Drop(ctx context.Context) error {
	fmt.Println("Dropping hotel collections")
	return s.collection.Drop(ctx)
}


func (s *MongoHotelStore) GetHotels(ctx context.Context, filter Map, pag *Pagination)([]*types.Hotel, error) {
	opts := options.FindOptions{}
	opts.SetSkip((pag.Page - 1) * pag.Limit)
	opts.SetLimit(pag.Limit)
	
	cur, err := s.collection.Find(ctx, filter, &opts)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := cur.All(ctx, &hotels); err != nil {
		return nil, err
	}

	return hotels, nil
}

func (s *MongoHotelStore) GetHotelById(ctx context.Context, id string) (*types.Hotel, error){
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var hotel types.Hotel
	if err := s.collection.FindOne(ctx,bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (s *MongoHotelStore) CreateHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter Map, params types.UpdateHotelParams) error {
	oid, err := primitive.ObjectIDFromHex(filter["_id"].(string))
	if err != nil {
		return  err
	}

	filter["_id"] = oid
	p := params.ToBSON()

	data := bson.M{}

	if _, ok := p["rooms"]; ok {
		for key, val := range p["rooms"].(bson.M) {
			data[key] = bson.M{"rooms": val}
		}
		delete(p, "rooms")
	}

	data["$set"] = p

	_, err = s.collection.UpdateOne(ctx, filter, data)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *MongoHotelStore) DeleteHotel(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return  err
	}
	
	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}