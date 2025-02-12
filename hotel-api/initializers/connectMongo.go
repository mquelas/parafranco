package initializers

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectMongo() {
	// Aquí debes colocar la URL de conexión de tu MongoDB.
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Conectando a MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Verificando la conexión
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Asignando la base de datos a la variable DB
	DB = client.Database("hotel_reservation")
}
