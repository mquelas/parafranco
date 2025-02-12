package main

import (
	"fmt"
	"os"
	"user-reservation-api/initializers"
	"user-reservation-api/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	// Verificar el valor de SECRET
	fmt.Println("SECRET:", os.Getenv("SECRET"))
	r := gin.Default()

	// Configuración CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"}, // Cambia esto por el origen correcto de tu frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rutas y controladores de usuarios y reservas
	routes.SetupUserRoutes(r)

	r.Run()
}
