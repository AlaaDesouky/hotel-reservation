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

const userCollection = "users"

type UserStore interface {
	Dropper

	GetUsers(context.Context) ([]*types.User, error)
	GetUserById(context.Context, string) (*types.User, error)
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(context.Context, Map, types.UpdateUserParams) error
	DeleteUser(context.Context, string) error
}

type MongoUserStore struct {
	client *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	db_name := os.Getenv(DB_NAME)
	return &MongoUserStore{
		client: client,
		collection: client.Database(db_name).Collection(userCollection),
	}
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("Dropping user collections")
	return s.collection.Drop(ctx)
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) GetUserById(ctx context.Context, id string) (*types.User, error){
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := s.collection.FindOne(ctx,bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}


func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error){
	var user types.User
	if err := s.collection.FindOne(ctx,bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter Map, params types.UpdateUserParams) error {
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

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
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