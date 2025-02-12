package main

import (
	"fmt"
	"hotel-api/initializers"
	"hotel-api/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Importar godotenv
)

func init() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file") // Si falla, detienes la ejecución
	}

	// Conectar a la base de datos
	initializers.ConnectMongo()

	// Conectar a RabbitMQ
	err := initializers.ConnectRabbitMQ() // Llamamos a la función para inicializar RabbitMQ
	if err != nil {
		panic("Failed to connect to RabbitMQ") // En caso de error, paramos la aplicación
	}
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

	// Llamar al archivo de rutas para registrar las rutas de los hoteles y amenidades
	routes.SetupHotelRoutes(r)
	routes.SetupAmenityRoutes(r)

	// Iniciar el servidor
	r.Run(":8080")
}
