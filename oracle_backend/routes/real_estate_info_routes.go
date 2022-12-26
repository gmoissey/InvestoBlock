package routes

import (
	"oracle_backend/controllers"

	"github.com/gin-gonic/gin"
)

func RealEstateInfoRoutes(router *gin.Engine) {
	router.GET("/real_estate_info", controllers.GetRealEstateByID())
	router.POST("/real_estate_info", controllers.CreateRealEstate())
}
