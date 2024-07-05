package db

import (
	"context"
	"fmt"
	"hotel-reservation/types"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingCollection = "bookings"

type BookingStore interface {
	Dropper

	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingById(context.Context, string) (*types.Booking, error)
	CreateBooking(context.Context, *types.Booking) (*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M)  error
}

type MongoBookingStore struct {
	client *mongo.Client
	collection *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	db_name := os.Getenv(DB_NAME)
	return &MongoBookingStore{
		client: client,
		collection: client.Database(db_name).Collection(bookingCollection),
	}
}

func (s *MongoBookingStore) Drop(ctx context.Context) error {
	fmt.Println("Dropping booking collections")
	return s.collection.Drop(ctx)
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error){
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var booking types.Booking
	if err := s.collection.FindOne(ctx,bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *MongoBookingStore) CreateBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}

	booking.ID = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, params bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return  err
	}

	data := bson.M{"$set": params}

	_, err = s.collection.UpdateByID(ctx, oid, data)
	if err != nil {
		return err
	}
	
	return nil
}