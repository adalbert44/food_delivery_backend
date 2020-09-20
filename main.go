package main

import (
	"food_delivery_backend/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(cors.Default())

	router.GET("/types", controllers.GetType)
	router.POST("/types", controllers.CreateType)
	router.PUT("/types", controllers.UpdateType)
	router.DELETE("/types", controllers.DeleteType)

	router.GET("/restaurants", controllers.GetRestaurant)
	router.POST("/restaurants", controllers.CreateRestaurant)
	router.PUT("/restaurants", controllers.UpdateRestaurant)
	router.DELETE("/restaurants", controllers.DeleteRestaurant)

	router.GET("/meals", controllers.GetMeal)
	router.POST("/meals", controllers.CreateMeal)
	router.PUT("/meals", controllers.UpdateMeal)
	router.DELETE("/meals", controllers.DeleteMeal)

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
