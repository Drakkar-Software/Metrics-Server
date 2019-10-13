package database

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

var dBName = getDBName()

// DB database structure
type DB struct {
	Client         *mongo.Client
	Collection     *mongo.Collection
	UserCollection *mongo.Collection
}

// Database is the database connexion instance
var Database = DB{}

// Init initializes the database connection
func Init() error {
	return Database.initialize()
}

// initialize connects to database and checks connection
func (db *DB) initialize() error {
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
	db.UserCollection = db.GetUserCollection()
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

func getDBName() string {
	urlData, err := url.Parse(DBURI())
	if err != nil {
		panic(err)
	}
	return urlData.Path[1:]
}

// DBURI returns the database URI
func DBURI() string {
	dburl := os.Getenv("MONGODB_URI")

	if dburl == "" {
		panic("Database url not found in environment variables. Please set MONGODB_URI env variable.")
	}

	return dburl
}

// GetCollection returns the metrics collection
func (db *DB) GetCollection() *mongo.Collection {
	return db.Client.Database(dBName).Collection("metrics")
}

// GetUserCollection returns the metrics collection
func (db *DB) GetUserCollection() *mongo.Collection {
	return db.Client.Database(dBName).Collection("users")
}
