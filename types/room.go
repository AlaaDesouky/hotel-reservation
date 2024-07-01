package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	minRoomPrice = 1.0
)

type UpdateRoomParams struct {
	Size    string  `json:"size"`
	Seaside *bool   `json:"seaside"`
	Price   float64 `json:"price"`
}

func (p UpdateRoomParams) ToBSON() bson.M {
	m := bson.M{}

	if len(p.Size) > 0 {
		m["size"] = p.Size
	}

	if p.Seaside != nil {
		m["seaside"] = *p.Seaside
	}

	if p.Price >= minRoomPrice {
		m["price"] = p.Price
	}

	return m
}

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Size    string             `bson:"size" json:"size"`
	Seaside bool               `bson:"seaside" json:"seaside"`
	Price   float64            `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
