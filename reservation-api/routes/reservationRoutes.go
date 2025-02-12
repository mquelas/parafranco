package routes

import (
	"reservation-api/controllers"

	"github.com/gin-gonic/gin"
)

// SetupReservationRoutes define las rutas para las operaciones de reserva
func SetupReservationRoutes(r *gin.Engine) {
	reservationGroup := r.Group("/reservations")
	{
		// Ruta para crear una nueva reserva
		reservationGroup.POST("/create", controllers.CreateReservation)

		// Ruta para obtener todas las reservas
		reservationGroup.GET("/all", controllers.GetAllReservations)

		// Ruta para obtener las reservas de un usuario por su ID
		reservationGroup.GET("/user/:userID", controllers.GetReservationsByUser)

		// Ruta para cancelar una reserva
		reservationGroup.DELETE("/cancel/:reservationID", controllers.CancelReservation)
	}
}
