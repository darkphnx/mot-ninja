package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

// Database wraps a mongo connection and adds convenience features
type Database struct {
	*mongo.Database
}

// InitDB establishes a database connection to Mongo
func InitDB(connectionString string) (*Database, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	db, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	database := Database{
		db.Database("vehicle-manager"),
	}

	return &database, nil
}
