package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Amenity representa una estructura de datos para las comodidades (amenities) de un hotel
type Amenity struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
}
