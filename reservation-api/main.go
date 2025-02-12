package main

import (
	"reservation-api/initializers"
	"reservation-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables() // Cargar variables de entorno
	initializers.ConnectToDb()      // Conectar a la base de datos
	initializers.SyncDatabase()     // Sincronizar la base de datos
}

func main() {
	r := gin.Default()

	// Configuración CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"}, // Cambia esto por el origen correcto de tu frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rutas y controladores de reservas
	routes.SetupReservationRoutes(r)

	// Ejecutar el servidor
	r.Run() // El puerto lo define desde el .env
}
