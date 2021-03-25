package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id"`
	Email          string             `bson:"email"`
	HashedPassword string             `bson:"hashed_password"`
	CreatedAt      time.Time          `bson:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

// CreateUser writes a new user to the database
func CreateUser(db *Database, user *User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := userCollection(db).InsertOne(ctx, user)
	return err
}

func GetUser(db *Database, email string) (*User, error) {
	var user User

	query := bson.M{
		"email": email,
	}

	err := userCollection(db).FindOne(ctx, query).Decode(&user)

	return &user, err
}

func userCollection(db *Database) *mongo.Collection {
	return db.Collection("users")
}
