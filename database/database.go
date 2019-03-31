package database

import (
	"context"
	"fmt"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

// DB database structure
type DB struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

// Initialize connects to database and checks connection
func (db *DB) Initialize() error {
	client, err := mongo.Connect(context.Background(), DBURI())
	if err != nil {
		return err
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return err
	}
	fmt.Println("Connected to database")
	db.Client = client
	db.Collection = db.GetCollection()
	return err
}

// Close from the database
func (db *DB) Close() error {
	err := db.Client.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	fmt.Println("Connection to database closed.")
	return err
}

// DBURI returns the database URI
func DBURI() string {
	dburl := os.Getenv("MONGODB_URI")

	if dburl == "" {
		dburl = "localhost"
	}

	return dburl
}

// GetCollection returns the metrics collection
func (db *DB) GetCollection() *mongo.Collection {
	return db.Client.Database("heroku_mq79r21v").Collection("metrics")
}
