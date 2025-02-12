package services

import (
	"context"
	"hotel-api/initializers" // Asegúrate de importar el paquete de inicialización
	"hotel-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Crear una amenidad
func CreateAmenity(amenityDto models.Amenity) (models.Amenity, error) {
	amenityDto.ID = primitive.NewObjectID()
	collection := initializers.DB.Collection("amenities") // Usa DB del archivo connectMongo

	_, err := collection.InsertOne(context.Background(), amenityDto)
	if err != nil {
		return models.Amenity{}, err
	}

	return amenityDto, nil
}

// Obtener una amenidad por ID
func GetAmenity(id primitive.ObjectID) (models.Amenity, error) {
	var amenity models.Amenity
	collection := initializers.DB.Collection("amenities")
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&amenity)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Amenity{}, err
		}
		return models.Amenity{}, err
	}
	return amenity, nil
}

func GetAmenities() ([]models.Amenity, error) {
	var amenities []models.Amenity
	collection := initializers.DB.Collection("amenities")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var amenity models.Amenity
		if err := cursor.Decode(&amenity); err != nil {
			return nil, err
		}
		amenities = append(amenities, amenity)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return amenities, nil
}

// Actualizar una amenidad
func UpdateAmenity(id primitive.ObjectID, amenityDto models.Amenity) (models.Amenity, error) {
	collection := initializers.DB.Collection("amenities")

	update := bson.M{
		"$set": amenityDto,
	}

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil {
		return models.Amenity{}, err
	}

	return amenityDto, nil
}

// Eliminar una amenidad
func DeleteAmenity(id primitive.ObjectID) error {
	collection := initializers.DB.Collection("amenities")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

// CheckDuplicate verifica si ya existe una amenidad con el mismo nombre
func CheckDuplicate(name string) (bool, error) {
	collection := initializers.DB.Collection("amenities")

	var result models.Amenity
	err := collection.FindOne(context.Background(), bson.M{"name": name}).Decode(&result)

	// Si no se encuentra ningún documento, significa que no hay duplicados
	if err == mongo.ErrNoDocuments {
		return false, nil
	}

	// Si encontramos un error distinto a "no documents", lo devolvemos
	if err != nil {
		return false, err
	}

	// Si llegamos aquí, es porque encontramos una amenidad con ese nombre
	return true, nil
}
