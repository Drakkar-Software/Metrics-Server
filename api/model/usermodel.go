package model

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Users access levels
const (
	FullAccess = 100
)

// User stores usage info about a specific user identified by Api-Key
type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	APIKey      string             `json:"apikey,omitempty" bson:"apiKey"`
	AccessLevel int8               `json:"accesslevel,omitempty" bson:"accessLevel"`
}
