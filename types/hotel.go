package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	minHotelNameLen = 2
	minHotelLocationLen = 2
	minHotelRating = 0
)

type UpdateHotelParams struct {
	Name string `json:"name"`
	Location string `json:"location"`
	Rating int `json:"rating"`
	Rooms bson.M `json:"rooms"`
}

func (p UpdateHotelParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.Name) >= minHotelNameLen {
		m["name"] = p.Name
	}
	if len(p.Location) >= minHotelLocationLen {
		m["location"] = p.Location
	}
	if p.Rating >= minHotelRating {
		m["rating"] = p.Rating
	}
	if len(p.Rooms) > 0 {
		m["rooms"] = p.Rooms
	}

	return m
}

type Hotel struct {
	ID  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `bson:"name" json:"name"`
	Location string `bson:"location" json:"location"`
	Rooms []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating int `bson:"rating" json:"rating"`
}
