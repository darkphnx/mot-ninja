package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Vehicle is a model of a vehicle inclusive of history that can be written to the database
type Vehicle struct {
	ID                 primitive.ObjectID `bson:"_id"`
	RegistrationNumber string             `bson:"registration_number"`
	Manufacturer       string             `bson:"manufacturer"`
	Model              string             `bson:"model"`
	MotDue             time.Time          `bson:"mot_due"`
	VEDDue             time.Time          `bson:"ved_due"`
}

// CreateVehicle writes a Vehicle struct to the database
func CreateVehicle(vehicle *Vehicle) error {
	_, err := collections.Vehicles.InsertOne(ctx, vehicle)
	return err
}
