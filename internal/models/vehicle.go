package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	MOTHistory         []MOTTest          `bson:"mot_history"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
}

// MOTTest that can be written to database
type MOTTest struct {
	TestNumber     int              `bson:"test_number"`
	Passed         bool             `bson:"passed"`
	CompletedDate  time.Time        `bson:"completed_date"`
	RfrAndComments []RfrAndComments `bson:"rfr_and_comments"`
}

// RfrAndComments contains the reasons for failure in a MOT
type RfrAndComments struct {
	Comment string `bson:"comment"`
	Type    string `bson:"type"`
}

// CreateVehicle writes a Vehicle struct to the database
func CreateVehicle(vehicle *Vehicle) error {
	vehicle.ID = primitive.NewObjectID()
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()

	_, err := collections.Vehicles.InsertOne(ctx, vehicle)
	return err
}

// GetVehicles fetches all vehicles
func GetVehicles() ([]*Vehicle, error) {
	var vehicles []*Vehicle

	cur, err := collections.Vehicles.Find(ctx, bson.D{{}})
	if err != nil {
		return vehicles, err
	}

	for cur.Next(ctx) {
		var vehicle Vehicle

		err := cur.Decode(&vehicle)
		if err != nil {
			return vehicles, err
		}

		vehicles = append(vehicles, &vehicle)
	}

	if err = cur.Err(); err != nil {
		return vehicles, err
	}

	cur.Close(ctx)

	return vehicles, nil
}
