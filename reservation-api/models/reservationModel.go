package models

import "time"

type Reservation struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	UserID     uint      `json:"userId"`
	HotelID    string    `json:"hotelId"`
	FechaDesde time.Time `json:"fechaDesde"`
	FechaHasta time.Time `json:"fechaHasta"`
}
