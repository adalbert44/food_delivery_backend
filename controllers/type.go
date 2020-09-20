package controllers

import (
	"fmt"
	"food_delivery_backend/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetType(c *gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	rows, err := conn.Query("SELECT * FROM FoodDelivery.type ORDER BY id")
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	types := make([]models.Type, 0)
	for rows.Next() {
		curType := models.Type {}
		err := rows.Scan(&curType.Id, &curType.Name)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}
		types = append(types, curType)
	}
	c.JSON(200, types)
}

func CreateType(c *gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	type Body struct {
		NewMealType models.Type `json:"newMealType"`
	}
	b := Body{}
	err = c.BindJSON(&b)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	_, err = conn.Exec("INSERT INTO FoodDelivery.type(name) VALUES (?)", b.NewMealType.Name)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	c.String(200, "OK")
}

func UpdateType(c *gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	type Body struct {
		EditMealType models.Type `json:"editMealType"`
	}
	b := Body{}
	err = c.BindJSON(&b)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	_, err = conn.Exec("UPDATE FoodDelivery.type SET name=? WHERE id=?", b.EditMealType.Name, b.EditMealType.Id)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	c.String(200, "OK")
}

func DeleteType(c *gin.Context) {
	conn, err := getDBConnection()
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	deleteMealTypeId, err := strconv.Atoi(c.Query("deleteMealTypeId"))
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	_, err = conn.Exec("DELETE FROM FoodDelivery.type WHERE id=?", deleteMealTypeId)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}
	_, err = conn.Exec("DELETE FROM FoodDelivery.meals WHERE typeid=?", deleteMealTypeId)
	if err != nil {
		c.String(404, fmt.Sprintf("%v", err))
		return
	}

	c.String(200, "OK")
}
