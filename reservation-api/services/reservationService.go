// reservationService.go

package services

import (
	"fmt"
	"net/http"
	"reservation-api/dto"
	"reservation-api/initializers"
	"reservation-api/models"
)

// CheckUserExists verifica si un usuario existe en la user-api
func CheckUserExists(userID uint) (bool, error) {
	url := fmt.Sprintf("http://localhost:3000/users/checkExistence/%d", userID) // Asegurar que la URL es correcta
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("error contacting user API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("unexpected response from user API: %d", resp.StatusCode)
}

// CheckHotelExists verifica si un hotel existe en la hotel-api
func CheckHotelExists(hotelID string) (bool, error) {
	url := fmt.Sprintf("http://localhost:8080/hotels/getHotel/%s", hotelID)
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("error contacting hotel API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else if resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	return false, fmt.Errorf("unexpected response from hotel API: %d", resp.StatusCode)
}

// CreateReservation crea una nueva reserva
func CreateReservation(reservationDto dto.ReservationDTO) (*models.Reservation, error) {
	// Verificar existencia del usuario
	userExists, err := CheckUserExists(reservationDto.UserID)
	if err != nil {
		return nil, fmt.Errorf("error checking user existence: %v", err)
	}
	if !userExists {
		return nil, fmt.Errorf("user does not exist")
	}

	// Verificar existencia del hotel
	hotelExists, err := CheckHotelExists(reservationDto.HotelID)
	if err != nil {
		return nil, fmt.Errorf("error checking hotel existence: %v", err)
	}
	if !hotelExists {
		return nil, fmt.Errorf("hotel does not exist")
	}

	// Si ambas verificaciones son exitosas, proceder con la reserva
	reservation := models.Reservation{
		UserID:     reservationDto.UserID,
		HotelID:    reservationDto.HotelID,
		FechaDesde: reservationDto.FechaDesde,
		FechaHasta: reservationDto.FechaHasta,
	}

	if err := initializers.DB.Create(&reservation).Error; err != nil {
		return nil, fmt.Errorf("failed to create reservation: %v", err)
	}

	return &reservation, nil
}

// GetAllReservations obtiene todas las reservas
func GetAllReservations() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := initializers.DB.Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetReservationsByUser obtiene todas las reservas de un usuario por su ID
func GetReservationsByUser(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	result := initializers.DB.Where("user_id = ?", userID).Find(&reservations)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch reservations for user %d: %v", userID, result.Error)
	}

	return reservations, nil
}

// CancelReservation cancela una reserva por su ID
func CancelReservation(reservationID uint) error {
	var reservation models.Reservation

	// Verificar si la reserva existe
	result := initializers.DB.Where("id = ?", reservationID).First(&reservation)
	if result.Error != nil {
		return fmt.Errorf("reservation not found")
	}

	// Eliminar la reserva
	result = initializers.DB.Delete(&reservation)
	if result.Error != nil {
		return fmt.Errorf("failed to delete reservation")
	}

	return nil
}
