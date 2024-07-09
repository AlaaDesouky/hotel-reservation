package fixtures

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, firstName, lastName string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email: fmt.Sprintf("%s@%s.com", firstName, lastName),
		FirstName: firstName,
		LastName: lastName,
		Password: fmt.Sprintf("%s_%s1234", firstName, lastName),
	})

	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin

	insertedUser, err := store.User.CreateUser(context.TODO(), user)

	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}


func AddHotel(store *db.Store, name, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []primitive.ObjectID{}
	}

	hotel:= types.Hotel{
		Name: name,
		Location: location,
		Rooms: roomIDS,
		Rating: rating,
	} 

	insertedHotel, err := store.Hotel.CreateHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := types.Room{
		Size: size,
		Seaside: seaside,
		Price: price,
		HotelID: hotelID,
	}

	insertedRoom, err := store.Room.CreateRoom(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(store *db.Store, uid, rid primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.CreateBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}