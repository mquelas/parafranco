package routes

import (
	"user-reservation-api/controllers"
	"user-reservation-api/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes define las rutas relacionadas con usuarios
func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", controllers.SignUp)                               // Registro de usuario
		userGroup.POST("/login", controllers.Login)                                   // Inicio de sesión
		userGroup.GET("/validate", middleware.RequireAuth, controllers.Validate)      // Validar sesión
		userGroup.GET("/current", middleware.RequireAuth, controllers.GetCurrentUser) // Obtener usuario actual
		userGroup.POST("/logout", controllers.Logout)                                 // Cerrar sesión
		// Nueva ruta para verificar si el usuario existe
		userGroup.GET("/checkExistence/:userID", controllers.CheckUserExistence) // Verificar existencia de usuario
	}
}
