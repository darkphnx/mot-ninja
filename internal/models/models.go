package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()
var database *mongo.Database

// InitDB establishes a database connection to Mongo
func InitDB(connectionString string) error {
	clientOptions := options.Client().ApplyURI(connectionString)
	db, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = db.Ping(ctx, nil)
	if err != nil {
		return err
	}

	database = db.Database("vehicle-manager")

	return nil
}
