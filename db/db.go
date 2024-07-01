package db

import "context"

const DB_NAME = "MONGO_DB_NAME"

type Map map[string]any
type Pagination struct {
	Limit int64
	Page int64
}

type Dropper interface {
	Drop(context.Context) error
}

type Store struct {
	User UserStore
	Hotel HotelStore
	Room RoomStore
}