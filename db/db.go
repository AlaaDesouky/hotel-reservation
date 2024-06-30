package db

import "context"

const DB_NAME = "MONGO_DB_NAME"


type Dropper interface {
	Drop(context.Context) error
}

type Store struct {
	User UserStore
}