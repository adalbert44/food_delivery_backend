package controllers

import (
	"fmt"
	"food_delivery_backend/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetMeal(c *gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	rows, err := conn.Query("SELECT * FROM FoodDelivery.meals ORDER BY id")
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	meals := make([]models.Meal, 0)
	for rows.Next() {
		curMeal := models.Meal {}
		err := rows.Scan(&curMeal.Id, &curMeal.Name, &curMeal.Description, &curMeal.PhotoURL, &curMeal.Price, &curMeal.TypeId)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		meals = append(meals, curMeal)
	}

	mealRequests := make([]models.MealRequest, 0)
	for _, curMeal :=  range meals {
		typeRows, err := conn.Query("SELECT * FROM FoodDelivery.type WHERE id = ? LIMIT 1", curMeal.TypeId)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}
		if !typeRows.Next() {
			c.String(404, fmt.Sprintf("Failed to find meal type with id %d", curMeal.TypeId))
			return
		}
		curType := models.Type{}
		err = typeRows.Scan(&curType.Id, &curType.Name)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		curMealRequest := models.MealRequest{
			Id:curMeal.Id,
			Name :curMeal.Name,
			Description:curMeal.Description,
			PhotoURL: curMeal.PhotoURL,
			Price:curMeal.Price,
			Type: curType,

		}

		mealRequests = append(mealRequests, curMealRequest)
	}

	c.JSON(200, mealRequests)
}

func CreateMeal(c * gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	type Body struct {
		NewMeal models.MealRequest `json:"newMeal"`
	}
	b := Body{}
	err = c.BindJSON(&b)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	rows, err := conn.Query("SELECT * FROM FoodDelivery.type WHERE id = ?", b.NewMeal.Type.Id)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}
	if !rows.Next() {
		c.String(404, fmt.Sprintf("meal type with id %d doesn't exist", b.NewMeal.Type.Id))
		return
	}

	_, err = conn.Exec("INSERT INTO FoodDelivery.meals(name, description, photourl, price, typeid) VALUES (?, ?, ?, ?, ?)",
		b.NewMeal.Name,
		b.NewMeal.Description,
		b.NewMeal.PhotoURL,
		b.NewMeal.Price,
		b.NewMeal.Type.Id,
	)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	c.String(200, "OK")
}

func UpdateMeal(c * gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	type Body struct {
		EditMeal models.MealRequest `json:"editMeal"`
	}
	b := Body{}
	err = c.BindJSON(&b)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	rows, err := conn.Query("SELECT * FROM FoodDelivery.type WHERE id = ?", b.EditMeal.Type.Id)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}
	if !rows.Next() {
		c.String(404, fmt.Sprintf("meal type with id %d doesn't exist", b.EditMeal.Type.Id))
		return
	}

	_, err = conn.Exec("UPDATE FoodDelivery.meals SET name=?, description=?, photourl=?, price=?, typeid=? WHERE id=?",
		b.EditMeal.Name,
		b.EditMeal.Description,
		b.EditMeal.PhotoURL,
		b.EditMeal.Price,
		b.EditMeal.Type.Id,
		b.EditMeal.Id,
	)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	c.String(200, "OK")
}

func DeleteMeal(c * gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	deleteMealId, err := strconv.Atoi(c.Query("deleteMealId"))
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	_, err = conn.Exec("DELETE FROM FoodDelivery.meals WHERE id=?", deleteMealId)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	c.String(200, "OK")
}