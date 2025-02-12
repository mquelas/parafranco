package controllers

import (
	"fmt"
	"hotel-api/initializers"
	"hotel-api/models"
	"hotel-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelController struct{}

// Crear un hotel
func (ctrl *HotelController) CreateHotel(c *gin.Context) {
	var hotelDto models.Hotel
	if err := c.ShouldBindJSON(&hotelDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el hotel ya existe con el mismo nombre y dirección
	duplicate, err := services.CheckDuplicateHotel(hotelDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for duplicate hotel"})
		return
	}

	if duplicate {
		c.JSON(http.StatusConflict, gin.H{"error": "Hotel with the same name and address already exists"})
		return
	}

	// Si no hay duplicado, crear el hotel
	hotel, err := services.CreateHotel(hotelDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hotel"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": hotel.ID.Hex()})
}

// Obtener un hotel por ID
func (ctrl *HotelController) GetHotel(c *gin.Context) {
	id := c.Param("id")

	// Convertir el ID de string a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	hotel, err := services.GetHotel(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

// Obtener todos los hoteles
func (ctrl *HotelController) GetHotels(c *gin.Context) {
	hotels, err := services.GetHotels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch hotels"})
		return
	}

	c.JSON(http.StatusOK, hotels)
}

// Actualizar un hotel
func (ctrl *HotelController) UpdateHotel(c *gin.Context) {
	id := c.Param("id")
	var hotelDto models.Hotel

	// Convertir el ID de string a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	if err := c.ShouldBindJSON(&hotelDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el hotel con el mismo nombre y dirección ya existe, pero excluyendo este hotel
	duplicate, err := services.CheckDuplicateHotel(hotelDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for duplicate hotel"})
		return
	}

	if duplicate {
		c.JSON(http.StatusConflict, gin.H{"error": "Hotel with the same name and address already exists"})
		return
	}

	// Si no hay duplicado, actualizar el hotel
	hotel, err := services.UpdateHotel(objectID, hotelDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hotel"})
		return
	}

	c.JSON(http.StatusOK, hotel)
}

// Eliminar un hotel
func (ctrl *HotelController) DeleteHotel(c *gin.Context) {
	id := c.Param("id")

	// Convertir el ID de string a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hotel ID"})
		return
	}

	err = services.DeleteHotel(objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete hotel"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Hotel deleted successfully"})
}

func CheckHotelExistence(c *gin.Context) {
	hotelID := c.Param("hotelID")

	var hotel models.Hotel
	// Buscar el hotel en la base de datos usando MongoDB
	if err := initializers.DB.Collection("hotels").FindOne(c, map[string]interface{}{"hotelID": hotelID}).Decode(&hotel); err != nil {
		// Si el hotel no existe, devolver 404
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Hotel with ID %s does not exist", hotelID)})
		return
	}

	// Si el hotel existe, devolver 200 OK
	c.JSON(http.StatusOK, gin.H{"message": "Hotel exists"})
}
