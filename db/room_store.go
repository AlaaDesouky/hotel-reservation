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

const roomCollection = "rooms"

type RoomStore interface {
	Dropper

	GetRooms(context.Context, Map) ([]*types.Room, error)
	GetRoomById(context.Context, string) (*types.Room, error)
	CreateRoom(context.Context, *types.Room) (*types.Room, error)
	UpdateRoom(context.Context, Map, types.UpdateRoomParams) error
	DeleteRoom(context.Context, string) error
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	db_name := os.Getenv(DB_NAME)
	return &MongoRoomStore{
		client: client,
		collection: client.Database(db_name).Collection(roomCollection),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) Drop(ctx context.Context) error {
	fmt.Println("Dropping room collections")
	return s.collection.Drop(ctx)
}


func (s *MongoRoomStore) GetRooms(ctx context.Context, filter Map)([]*types.Room, error) {
	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) GetRoomById(ctx context.Context, id string) (*types.Room, error){
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var room types.Room
	if err := s.collection.FindOne(ctx,bson.M{"_id": oid}).Decode(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (s *MongoRoomStore) CreateRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = res.InsertedID.(primitive.ObjectID)

	// Add room to hotel
	filter := Map{"_id": room.HotelID}
	update := types.UpdateHotelParams{
		Rooms: []primitive.ObjectID{room.ID},
	}

	if err := s.HotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *MongoRoomStore) UpdateRoom(ctx context.Context, filter Map, params types.UpdateRoomParams) error {
	oid, err := primitive.ObjectIDFromHex(filter["_id"].(string))
	if err != nil {
		return  err
	}

	filter["_id"] = oid
	data := bson.M{"$set": params.ToBSON()}

	_, err = s.collection.UpdateOne(ctx, filter, data)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *MongoRoomStore) DeleteRoom(ctx context.Context, id string) error {
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