package api

import (
	"context"
	"fmt"
	"hotel-reservation/db"
	"hotel-reservation/types"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

type BookRoomParams struct {
	FromDate time.Time `json:"fromDate"`
	TillDate time.Time `json:"tillDate"`
	NumPersons int `json:"numPersons"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("invalid date for booking a room")
	}
	return nil
}

func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), db.Map{})
	if err != nil {
		return ErrorResourceNotFound("rooms")
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	if err := params.validate(); err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return ErrorBadRequest()
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return ErrorInternalServer()
	}

	ok, err = h.isRoomAvailable(c.Context(), roomID, params)
	if err != nil {
		return err
	}

	if !ok {
		return c.Status(http.StatusBadRequest).JSON(GenericResponse{
			Type: "error",
			Msg: fmt.Sprintf("room %s already booked", c.Params("id")),
		})
	}

	booking := types.Booking{
		UserID: user.ID,
		RoomID: roomID,
		FromDate: params.FromDate,
		TillDate: params.TillDate,
		NumPersons: params.NumPersons,
	}

	insertedBooking, err := h.store.Booking.CreateBooking(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(insertedBooking)
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"formDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}

	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0

	return ok, nil
}