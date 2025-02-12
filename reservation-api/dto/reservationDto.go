package dto

import "time"

// DTO para la respuesta de la reserva completa
type ReservationDTO struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"userId"`
	HotelID    string    `json:"hotelId"`
	FechaDesde time.Time `json:"fechaDesde"`
	FechaHasta time.Time `json:"fechaHasta"`
}
