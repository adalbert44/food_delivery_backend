package controllers

import (
	"fmt"
	"food_delivery_backend/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetRestaurant(c *gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	rows, err := conn.Query("SELECT * FROM FoodDelivery.restaurants")
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	restaurants := make([]models.Restaurant, 0)
	for rows.Next() {
		curRestaurant := models.Restaurant{}
		err := rows.Scan(&curRestaurant.Id,
			&curRestaurant.Name,
			&curRestaurant.Description,
			&curRestaurant.Location,
			&curRestaurant.PhotoURL)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}
		restaurants = append(restaurants, curRestaurant)
	}
	c.JSON(200, restaurants)
}

func CreateRestaurant(c *gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	type Body struct {
		NewRestaurant models.Restaurant `json:"newRestaurant"`
	}
	b := Body{}
	err = c.BindJSON(&b)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	_, err = conn.Exec("INSERT INTO FoodDelivery.restaurants(name, description, location, photourl) VALUES (?, ?, ?, ?)",
		b.NewRestaurant.Name,
		b.NewRestaurant.Description,
		b.NewRestaurant.Location,
		b.NewRestaurant.PhotoURL,
	)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	c.String(200, "OK")
}

func UpdateRestaurant(c *gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	type Body struct {
		EditRestaurant models.Restaurant `json:"editRestaurant"`
	}
	b := Body{}
	err = c.BindJSON(&b)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	_, err = conn.Exec("UPDATE FoodDelivery.restaurants SET name=?, description=?, location=?, photourl=? WHERE id=?",
		b.EditRestaurant.Name,
		b.EditRestaurant.Description,
		b.EditRestaurant.Location,
		b.EditRestaurant.PhotoURL,
		b.EditRestaurant.Id,
	)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	c.String(200, "OK")
}

func DeleteRestaurant(c *gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	deleteRestaurantId, err := strconv.Atoi(c.Query("deleteRestaurantId"))
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	_, err = conn.Exec("DELETE FROM FoodDelivery.restaurants WHERE id=?", deleteRestaurantId)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	c.String(200, "OK")
}
