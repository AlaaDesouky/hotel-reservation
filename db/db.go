package db

const (
	DB_NAME      = "MONGO_DB_NAME"
	TEST_DB_NAME = "MONGO_TEST_DB_NAME"
)

type Store struct {
	User UserStore
}