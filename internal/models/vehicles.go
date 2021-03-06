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

var collectionName = "vehicles"

// CreateVehicle writes a Vehicle struct to the database
func CreateVehicle(db *Database, vehicle *Vehicle) error {
	vehicle.ID = primitive.NewObjectID()
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()

	_, err := db.Collection(collectionName).InsertOne(ctx, vehicle)
	return err
}

// GetAllVehicles fetches all vehicles
func GetAllVehicles(db *Database) ([]*Vehicle, error) {
	query := bson.M{}

	return getVehicles(db, query)
}

func getVehicles(db *Database, query bson.M) ([]*Vehicle, error) {
	var vehicles []*Vehicle

	cur, err := db.Collection(collectionName).Find(ctx, query)
	if err != nil {
		return vehicles, err
	}

	err = cur.All(ctx, &vehicles)
	if err != nil {
		return vehicles, err
	}

	return vehicles, nil
}
