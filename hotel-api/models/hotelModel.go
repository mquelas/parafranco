package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
    ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    Name      string             `json:"name" bson:"name"`
    Address   string             `json:"address" bson:"address"`
    City      string             `json:"city" bson:"city"`
    Country   string             `json:"country" bson:"country"`
    Amenities []string           `json:"amenities" bson:"amenities"`
    Photos    interface{}        `json:"photos" bson:"photos"` // Esto puede ser ajustado según el tipo de datos
}
