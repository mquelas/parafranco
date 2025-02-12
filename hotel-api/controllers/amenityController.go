package controllers

import (
	"hotel-api/dtos"
	"hotel-api/models"
	"hotel-api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AmenityController struct{}

// crear amenity
func (ctrl *AmenityController) CreateAmenity(c *gin.Context) {
	var amenityDto dtos.AmenityDto
	if err := c.ShouldBindJSON(&amenityDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si ya existe una amenidad con el mismo nombre
	exists, err := services.CheckDuplicate(amenityDto.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for duplicate amenity"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Amenity with this name already exists"})
		return
	}

	amenity := models.Amenity{
		Name: amenityDto.Name,
	}

	createdAmenity, err := services.CreateAmenity(amenity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create amenity"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ID": createdAmenity.ID.Hex()})
}

// Obtener un amenity por ID
func (ctrl *AmenityController) GetAmenity(c *gin.Context) {
	id := c.Param("id")

	// Convertir el ID de string a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amenity ID"})
		return
	}

	amenity, err := services.GetAmenity(objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Amenity not found"})
		return
	}

	c.JSON(http.StatusOK, amenity)
}

// Obtener todos los amenities
func (ctrl *AmenityController) GetAmenities(c *gin.Context) {
	amenities, err := services.GetAmenities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch amenities"})
		return
	}

	c.JSON(http.StatusOK, amenities)
}

// Actualizar un amenity
func (ctrl *AmenityController) UpdateAmenity(c *gin.Context) {
	id := c.Param("id")
	var amenityDto dtos.AmenityDto

	// Convertir el ID de string a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amenity ID"})
		return
	}

	if err := c.ShouldBindJSON(&amenityDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar si ya existe una amenidad con el mismo nombre, excluyendo la actual
	exists, err := services.CheckDuplicate(amenityDto.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for duplicate amenity"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Amenity with this name already exists"})
		return
	}

	amenity := models.Amenity{
		ID:   objectID,
		Name: amenityDto.Name,
	}

	updatedAmenity, err := services.UpdateAmenity(objectID, amenity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update amenity"})
		return
	}

	c.JSON(http.StatusOK, updatedAmenity)
}

// Eliminar un amenity
func (ctrl *AmenityController) DeleteAmenity(c *gin.Context) {
	id := c.Param("id")

	// Convertir el ID de string a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amenity ID"})
		return
	}

	err = services.DeleteAmenity(objectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete amenity"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Amenity deleted successfully"})
}
