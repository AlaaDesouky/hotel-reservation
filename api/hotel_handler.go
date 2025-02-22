package api

import (
	"hotel-reservation/db"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrorBadRequest()
	}

	filter := db.Map{}

	if params.Rating > 0 {
		filter["Rating"] = params.Rating 
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrorResourceNotFound("hotels")
	}

	resp := ResourceResponse{
		Data: hotels,
		Results: len(hotels),
		Page: int(params.Page),
	}

	return c.JSON(resp)
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelById(c.Context(), id)
	if err != nil {
		return ErrorResourceNotFound("hotel")
	}

	return c.JSON(hotel)
}

func (h *HotelHandler) HandleGetHotelRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrorInvalidID()
	}

	filter := db.Map{"hotelID": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrorResourceNotFound("hotel_rooms")
	}

	return c.JSON(rooms)
}