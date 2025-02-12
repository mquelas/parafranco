package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reservation-api/dto"
	"reservation-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Verificar usuario en user-api
func GetAuthenticatedUser(c *gin.Context) (uint, error) {
	// Obtener el token desde la cookie
	token, err := c.Cookie("Authorization")
	if err != nil {
		return 0, fmt.Errorf("No autorizado: Token no encontrado")
	}

	// Hacer una solicitud HTTP a user-api para verificar el token
	url := "http://localhost:3000/users/me" // Endpoint en user-api para validar usuario
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Cookie", "Authorization="+token) // Pasar el token en la cookie

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("No autorizado: No se pudo autenticar el usuario")
	}
	defer resp.Body.Close()

	// Leer la respuesta
	body, _ := ioutil.ReadAll(resp.Body)
	userID, err := ParseUserID(body)
	if err != nil {
		return 0, fmt.Errorf("Error al obtener ID de usuario")
	}

	return userID, nil
}

// Parsear el ID del usuario desde la respuesta de user-api
func ParseUserID(body []byte) (uint, error) {
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return 0, err
	}

	// Extraer el userID
	userIDFloat, ok := response["id"].(float64)
	if !ok {
		return 0, fmt.Errorf("ID de usuario inválido")
	}

	return uint(userIDFloat), nil
}

// Crear una reserva
func CreateReservation(c *gin.Context) {
	// Verificar usuario autenticado en user-api
	userID, err := GetAuthenticatedUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Vincular el JSON recibido a la estructura DTO
	var dto dto.ReservationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Asignar el UserID autenticado
	dto.UserID = userID

	// Llamar al servicio para crear la reserva (sin `c`)
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
