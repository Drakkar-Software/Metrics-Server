package model

// Top stores data about a ranking
type Top struct {
	Name  string `json:"name" bson:"_id"`
	Count int    `json:"count" bson:"count"`
}

const (
	AllTraders       = "AllTraders"
	SimulatedTraders = "SimulatedTraders"
	RealTraders      = "RealTraders"
)
