package api

import (
	"hotel-reservation/db"

	"github.com/gofiber/fiber/v2"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), db.Map{})
	if err != nil {
		return ErrorResourceNotFound("rooms")
	}

	return c.JSON(rooms)
}