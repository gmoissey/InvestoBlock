package controllers

import (
	"context"
	"net/http"
	"oracle_backend/database"
	"oracle_backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()
var realEstateInfoCollection = database.GetCollection(database.DB, "real_estate")

func CreateRealEstate() gin.HandlerFunc {
	return func(c *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//validate the request body
		var newRealEstate models.RealEstate
		if err := c.ShouldBindJSON(&newRealEstate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//use the validator library to validate required fields
		if err := validate.Struct(newRealEstate); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//insert the new real estate into the database
		result, err := realEstateInfoCollection.InsertOne(context, newRealEstate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//return the new real estate
		c.JSON(http.StatusOK, gin.H{
			"message": "Real estate created successfully",
			"result":  result,
		})
	}
}

func GetRealEstateByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var realEstate models.RealEstate
		idString := c.Param("id")
		id, _ := primitive.ObjectIDFromHex(idString)

		//find the real estate by id
		err := realEstateInfoCollection.FindOne(context, bson.M{"_id": id}).Decode(&realEstate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//return the real estate
		c.JSON(http.StatusOK, gin.H{
			"message": "Real estate found successfully",
			"result":  realEstate,
		})
	}
}

func UpdateRealEstateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//validate the request body
		var realEstate models.RealEstate
		if err := ctx.ShouldBindJSON(&realEstate); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//use the validator library to validate required fields
		if err := validate.Struct(realEstate); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//get the real estate id
		idString := ctx.Param("id")
		id, _ := primitive.ObjectIDFromHex(idString)

		update := bson.M{
			"market_price": realEstate.MarketPrice,
			"address":      realEstate.Address,
			"city":         realEstate.City,
			"state":        realEstate.State,
			"zip_code":     realEstate.ZipCode,
			"beds":         realEstate.Beds,
			"baths":        realEstate.Baths,
			"sqft":         realEstate.Sqft,
			"year_built":   realEstate.YearBuilt,
		}

		//update the real estate
		result, err := realEstateInfoCollection.UpdateOne(context, bson.M{"_id": id}, bson.M{"$set": update})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//return the updated real estate
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Real estate updated successfully",
			"result":  result,
		})
	}
}

func DeleteRealEstateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//get the real estate id
		idString := ctx.Param("id")
		id, _ := primitive.ObjectIDFromHex(idString)

		//delete the real estate
		result, err := realEstateInfoCollection.DeleteOne(context, bson.M{"_id": id}) //delete the real estate
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//return the deleted real estate
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Real estate deleted successfully",
			"result":  result,
		})
	}
}

func GetFilteredRealEstateInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//get the query params
		marketPriceMin := ctx.Query("market_price_min")
		marketPriceMax := ctx.Query("market_price_max")
		states := ctx.Query("states")
		cities := ctx.Query("cities")
		bedsMin := ctx.Query("beds_min")
		bedsMax := ctx.Query("beds_max")
		bathsMin := ctx.Query("baths_min")
		bathsMax := ctx.Query("baths_max")
		sqftMin := ctx.Query("sqft_min")
		sqftMax := ctx.Query("sqft_max")
		yearBuiltMin := ctx.Query("year_built_min")
		yearBuiltMax := ctx.Query("year_built_max")

		//create the filter
		filter := bson.M{}

		//add the market price filter
		if marketPriceMin != "" && marketPriceMax != "" {
			filter["market_price"] = bson.M{"$gte": 130000, "$lte": 230000}
		}

		//add the state filter
		if states != "" {
			filter["state"] = bson.M{"$in": strings.Split(states, ",")}
		}

		//add the city filter
		if cities != "" {
			filter["city"] = bson.M{"$in": strings.Split(cities, ",")}
		}

		//add the beds filter
		if bedsMin != "" && bedsMax != "" {
			filter["beds"] = bson.M{"$gte": bedsMin, "$lte": bedsMax}
		}

		//add the baths filter
		if bathsMin != "" && bathsMax != "" {
			filter["baths"] = bson.M{"$gte": bathsMin, "$lte": bathsMax}
		}

		//add the sqft filter
		if sqftMin != "" && sqftMax != "" {
			filter["sqft"] = bson.M{"$gte": sqftMin, "$lte": sqftMax}
		}

		//add the year built filter
		if yearBuiltMin != "" && yearBuiltMax != "" {
			filter["year_built"] = bson.M{"$gte": yearBuiltMin, "$lte": yearBuiltMax}
		}

		//get the real estate info
		cursor, err := realEstateInfoCollection.Find(context, filter)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//close the cursor
		defer cursor.Close(context)

		//create a slice to store the real estate info
		var realEstates []models.RealEstate

		//loop through the cursor
		for cursor.Next(context) {
			//create a real estate
			var realEstate models.RealEstate

			//decode the real estate
			cursor.Decode(&realEstate)

			//add the real estate to the slice
			realEstates = append(realEstates, realEstate)
		}

		//return the real estate info
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Real estate info retrieved successfully",
			"result":  realEstates,
		})
	}
}
