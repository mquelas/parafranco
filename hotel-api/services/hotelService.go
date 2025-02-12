package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hotel-api/initializers"
	"hotel-api/models"
	"log"
	"strings"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Validar si las amenidades existen
func validateAmenitiesExist(amenities []string) error {
	var invalidAmenities []string

	// Recorremos las amenidades y verificamos si existen en la base de datos
	for _, amenity := range amenities {
		var existingAmenity models.Amenity
		err := initializers.DB.Collection("amenities").FindOne(context.Background(), bson.M{"name": amenity}).Decode(&existingAmenity)
		if err != nil {
			invalidAmenities = append(invalidAmenities, amenity)
		}
	}

	if len(invalidAmenities) > 0 {
		// Devolvemos un error si alguna amenidad no existe
		return errors.New("the following amenities do not exist: " + strings.Join(invalidAmenities, ", "))
	}

	return nil
}

// Enviar un mensaje a RabbitMQ
func SendHotelCreationMessage(hotel models.Hotel) error {
	if initializers.RabbitMQChannel == nil {
		log.Println("RabbitMQ channel is not initialized")
		return fmt.Errorf("RabbitMQ channel is not initialized")
	}

	message := map[string]interface{}{
		"id":        hotel.ID.Hex(),
		"name":      hotel.Name,
		"address":   hotel.Address,
		"city":      hotel.City,
		"country":   hotel.Country,
		"amenities": hotel.Amenities,
	}

	// Convertir el mensaje a JSON
	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal hotel message: %s", err)
		return err
	}

	// Asegurarse de que la cola "hotel_created" esté declarada
	_, err = initializers.RabbitMQChannel.QueueDeclare(
		"hotel_created", // nombre de la cola
		true,            // durable
		false,           // auto-deleted
		false,           // exclusive
		false,           // no-wait
		nil,             // argumentos adicionales
	)
	if err != nil {
		log.Printf("Failed to declare queue: %s", err)
		return err
	}

	// Publicar el mensaje
	err = initializers.RabbitMQChannel.Publish(
		"",              // exchange
		"hotel_created", // routing key (nombre de la cola)
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %s", err)
		return err
	}

	log.Println("Message sent to RabbitMQ queue: hotel_created")
	return nil
}

// Crear un hotel
func CreateHotel(hotelDto models.Hotel) (models.Hotel, error) {
	// Verificar que las amenidades existan
	if err := validateAmenitiesExist(hotelDto.Amenities); err != nil {
		return models.Hotel{}, err
	}

	hotelDto.ID = primitive.NewObjectID()

	collection := initializers.DB.Collection("hotels")
	_, err := collection.InsertOne(context.Background(), hotelDto)
	if err != nil {
		return models.Hotel{}, err
	}

	// Enviar mensaje a RabbitMQ después de crear el hotel
	if err := SendHotelCreationMessage(hotelDto); err != nil {
		log.Printf("Failed to send hotel creation message: %s", err)
		return models.Hotel{}, err
	}

	return hotelDto, nil
}

// Obtener un hotel por ID
func GetHotel(id primitive.ObjectID) (models.Hotel, error) {
	var hotel models.Hotel
	collection := initializers.DB.Collection("hotels")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&hotel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Hotel{}, err
		}
		return models.Hotel{}, err
	}

	return hotel, nil
}

// Obtener todos los hoteles
func GetHotels() ([]models.Hotel, error) {
	var hotels []models.Hotel
	collection := initializers.DB.Collection("hotels")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var hotel models.Hotel
		if err := cursor.Decode(&hotel); err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return hotels, nil
}

// Actualizar un hotel
func UpdateHotel(id primitive.ObjectID, hotelDto models.Hotel) (models.Hotel, error) {
	// Verificar que las amenidades existan
	if err := validateAmenitiesExist(hotelDto.Amenities); err != nil {
		return models.Hotel{}, err
	}

	collection := initializers.DB.Collection("hotels")

	update := bson.M{
		"$set": hotelDto,
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return models.Hotel{}, err
	}

	return hotelDto, nil
}

// Eliminar un hotel
func DeleteHotel(id primitive.ObjectID) error {
	collection := initializers.DB.Collection("hotels")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

// CheckDuplicateHotel verifica si ya existe un hotel con el mismo nombre y dirección
func CheckDuplicateHotel(hotelDto models.Hotel) (bool, error) {
	collection := initializers.DB.Collection("hotels")

	// Buscar si hay un hotel con el mismo nombre y dirección
	var existingHotel models.Hotel
	filter := bson.M{
		"name":    hotelDto.Name,
		"address": hotelDto.Address,
	}
	err := collection.FindOne(context.Background(), filter).Decode(&existingHotel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No hay hotel duplicado
			return false, nil
		}
		return false, err // Otro error, posiblemente de conexión a la base de datos
	}

	// Si encontramos un hotel, devolvemos true
	return true, nil
}
