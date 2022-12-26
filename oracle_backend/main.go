package main

import (
	"oracle_backend/database"
	"oracle_backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	database.ConnectDB()
	routes.RealEstateInfoRoutes(router)
	router.Run("localhost:6000")
}
