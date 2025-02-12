package routes

import (
	"hotel-api/controllers"
	"hotel-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupAmenityRoutes(r *gin.Engine) {
	amenityController := controllers.AmenityController{}

	// Rutas protegidas con autenticación
	auth := r.Group("/")
	auth.Use(middleware.RequireAuth)

	// Rutas restringidas solo para administradores
	admin := auth.Group("/")
	admin.Use(middleware.RequireAdmin)

	// Solo administradores pueden crear, actualizar y eliminar amenities
	admin.POST("/createAmenity", amenityController.CreateAmenity)
	admin.PUT("/updateAmenity/:id", amenityController.UpdateAmenity)
	admin.DELETE("/deleteAmenity/:id", amenityController.DeleteAmenity)

	// Todos los usuarios pueden obtener amenities
	auth.GET("/getAmenityByID/:id", amenityController.GetAmenity)
	auth.GET("/getAllAmenities", amenityController.GetAmenities)
}
