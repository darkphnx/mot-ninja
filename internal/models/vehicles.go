package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Vehicle is a model of a vehicle inclusive of history that can be written to the database
type Vehicle struct {
	ID                 primitive.ObjectID `bson:"_id"`
	UserID             primitive.ObjectID `bson:"user_id"`
	RegistrationNumber string             `bson:"registration_number"`
	Manufacturer       string             `bson:"manufacturer"`
	Model              string             `bson:"model"`
	MotDue             time.Time          `bson:"mot_due"`
	VEDDue             time.Time          `bson:"ved_due"`
	MOTHistory         []MOTTest          `bson:"mot_history"`
	CreatedAt          time.Time          `bson:"created_at"`
	UpdatedAt          time.Time          `bson:"updated_at"`
	LastFetchedAt      time.Time          `bson:"last_fetched_at"`
}

// MOTTest that can be written to database
type MOTTest struct {
	TestNumber      int              `bson:"test_number"`
	Passed          bool             `bson:"passed"`
	CompletedDate   time.Time        `bson:"completed_date"`
	ExpiryDate      time.Time        `bson:"expiry_date"`
	OdometerReading string           `bson:"odometer_reading"`
	RfrAndComments  []RfrAndComments `bson:"rfr_and_comments"`
}

// RfrAndComments contains the reasons for failure in a MOT
type RfrAndComments struct {
	Comment string `bson:"comment"`
	Type    string `bson:"type"`
}

// CreateVehicle writes a Vehicle struct to the database
func CreateVehicle(db *Database, vehicle *Vehicle) error {
	vehicle.ID = primitive.NewObjectID()
	vehicle.CreatedAt = time.Now()
	vehicle.UpdatedAt = time.Now()

	_, err := vehicleCollection(db).InsertOne(ctx, vehicle)
	return err
}

func GetUserVehicle(db *Database, userID primitive.ObjectID, registrationNumber string) (*Vehicle, error) {
	var vehicle Vehicle

	query := bson.M{
		"registration_number": registrationNumber,
		"user_id":             userID,
	}

	err := vehicleCollection(db).FindOne(ctx, query).Decode(&vehicle)

	return &vehicle, err
}

// DeleteVehicle deletes a vehicle from the database
func DeleteVehicle(db *Database, vehicle *Vehicle) error {
	_, err := vehicleCollection(db).DeleteOne(ctx, bson.M{"_id": primitive.ObjectID(vehicle.ID)})

	return err
}

// GetUserVehicles fetches all vehicles for the given user ID
func GetUserVehicles(db *Database, userID primitive.ObjectID) ([]*Vehicle, error) {
	query := bson.M{
		"user_id": userID,
	}

	return getVehicles(db, query)
}

// GetVehiclesUpdatedBefore fetches any vehicle that has a LastRemotePull value less than timestamp
func GetVehiclesUpdatedBefore(db *Database, timestamp time.Time) ([]*Vehicle, error) {
	query := bson.M{
		"last_fetched_at": bson.M{"$lt": timestamp},
	}

	return getVehicles(db, query)
}

// UpdateVehicle replaces the existing vehicle with a brand new one
func UpdateVehicle(db *Database, existing *Vehicle, updated *Vehicle) error {
	updated.ID = existing.ID

	_, err := vehicleCollection(db).ReplaceOne(
		ctx,
		bson.M{"_id": primitive.ObjectID(existing.ID)},
		updated,
	)

	return err
}

func vehicleCollection(db *Database) *mongo.Collection {
	return db.Collection("vehicles")
}

func getVehicles(db *Database, query bson.M) ([]*Vehicle, error) {
	var vehicles []*Vehicle

	cur, err := vehicleCollection(db).Find(ctx, query)
	if err != nil {
		return vehicles, err
	}

	err = cur.All(ctx, &vehicles)
	if err != nil {
		return vehicles, err
	}

	return vehicles, nil
}
