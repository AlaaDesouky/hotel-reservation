package api

import (
	"hotel-reservation/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}

func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrorResourceNotFound("bookings")
	}

	return c.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return ErrorResourceNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return ErrorUnAuthorized()
	}

	if booking.UserID != user.ID {
		return ErrorUnAuthorized()
	}

	return c.JSON(booking)
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return ErrorResourceNotFound("booking")
	}

	user, err := getAuthUser(c)
	if err != nil {
		return ErrorUnAuthorized()
	}

	if booking.UserID != user.ID {
		return ErrorUnAuthorized()
	}

	if err := h.store.Booking.UpdateBooking(c.Context(), c.Params("id"), bson.M{"cancelled": true}); err != nil {
		return err
	}

	return c.JSON(GenericResponse{
		Type: "msg",
		Msg: "updated",
	})
}