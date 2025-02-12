package controllers

import (
	"net/http"
	"reservation-api/dto"
	"reservation-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateReservation crea una nueva reserva
func CreateReservation(c *gin.Context) {
	var dto dto.ReservationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	reservation, err := services.CreateReservation(dto)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"reservation": reservation})
}

// GetAllReservations obtiene todas las reservas
func GetAllReservations(c *gin.Context) {
	reservations, err := services.GetAllReservations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

// GetReservationsByUser obtiene las reservas de un usuario por su ID
func GetReservationsByUser(c *gin.Context) {
	userID := c.Param("userID")
	userIDInt, err := strconv.ParseUint(userID, 10, 32) // Convertir de string a uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	reservations, err := services.GetReservationsByUser(uint(userIDInt)) // Pasar como uint
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations for user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reservations": reservations})
}

// CancelReservation cancela una reserva por su ID
func CancelReservation(c *gin.Context) {
	reservationID := c.Param("reservationID")
	reservationIDInt, err := strconv.ParseUint(reservationID, 10, 32) // Convertir de string a uint
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
		return
	}

	err = services.CancelReservation(uint(reservationIDInt)) // Pasar como uint
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation canceled successfully"})
}
