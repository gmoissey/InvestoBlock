package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RealEstate struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MarketPrice float32            `json:"market_price,omitempty" bson:"market_price,omitempty"`
	Address     string             `json:"address,omitempty" bson:"address,omitempty"`
	ZipCode     string             `json:"zip_code,omitempty" bson:"zip_code,omitempty"`
	City        string             `json:"city,omitempty" bson:"city,omitempty"`
	State       string             `json:"state,omitempty" bson:"state,omitempty"`
	Beds        int                `json:"beds,omitempty" bson:"beds,omitempty"`
	Baths       int                `json:"baths,omitempty" bson:"baths,omitempty"`
	Sqft        int                `json:"sqft,omitempty" bson:"sqft,omitempty"`
	YearBuilt   int                `json:"year_built,omitempty" bson:"year_built,omitempty"`
}
