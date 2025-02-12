package dto

import "time"

type ReservationDTO struct {
	HotelID    string    `json:"hotelId" binding:"required"`
	FechaDesde time.Time `json:"fechaDesde" binding:"required"`
	FechaHasta time.Time `json:"fechaHasta" binding:"required"`
	UserID     uint      `json:"-"` // Evitar que el usuario lo pase manualmente
}
