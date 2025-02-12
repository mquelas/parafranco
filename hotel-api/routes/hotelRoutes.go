package routes

import (
	"hotel-api/controllers"
	"hotel-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupHotelRoutes(r *gin.Engine) {
	hotelController := &controllers.HotelController{}

	// Grupo de rutas protegidas con autenticación
	auth := r.Group("/hotels")
	auth.Use(middleware.RequireAuth) // Requiere autenticación

	// Grupo de rutas solo para administradores
	admin := auth.Group("")            // No es necesario repetir "/hotels"
	admin.Use(middleware.RequireAdmin) // Requiere rol admin

	// Solo administradores pueden crear, actualizar y eliminar hoteles
	admin.POST("/createHotel", hotelController.CreateHotel)
	admin.PUT("/updateHotel/:id", hotelController.UpdateHotel)
	admin.DELETE("/deleteHotel/:id", hotelController.DeleteHotel)

	// Todos los usuarios autenticados pueden ver hoteles
	auth.GET("/getHotels", hotelController.GetHotels)
	auth.GET("/getHotel/:id", hotelController.GetHotel)
	auth.GET("/check-existence/:hotelID", controllers.CheckHotelExistence)
}
