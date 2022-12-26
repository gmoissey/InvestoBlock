package controllers

import (
	"context"
	"net/http"
	"oracle_backend/database"
	"oracle_backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
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
		err := realEstateInfoCollection.FindOne(context, models.RealEstate{ID: id}).Decode(&realEstate)
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
