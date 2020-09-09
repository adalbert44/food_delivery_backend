package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Type struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Restaurant struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Location string `json:"location"`
	PhotoURL string `json:"photoUrl"`
}

type Meal struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	PhotoURL string `json:"photoUrl"`
	Price int `json:"price"`
	TypeId int `json:"type"`
}

func getDBConnection() (*sqlx.DB, error){
	connParams := strings.Join([]string{
		"parseTime=true",
		"interpolateParams=true",
		"timeout=10s",
		"collation_server=utf8_general_ci",
		"sql_select_limit=18446744073709551615",
		"compile_only=false",
		"enable_auto_profile=false",
		"sql_mode='STRICT_ALL_TABLES,ONLY_FULL_GROUP_BY'",
	}, "&")

	defaultConfigString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/information_schema?%s",
		os.Getenv("MEMSQL_USER"),
		os.Getenv("MEMSQL_PASSWORD"),
		os.Getenv("MEMSQL_HOST"),
		os.Getenv("MEMSQL_PORT"),
		connParams,
	)

	return sqlx.Open("mysql", defaultConfigString)
}

type MealRequest struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	PhotoURL string `json:"photoUrl"`
	Price int `json:"price"`
	Type Type `json:"type"`
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(cors.Default())

	router.GET("/types", func(c *gin.Context) {
		conn, err := getDBConnection()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		rows, err := conn.Query("SELECT * FROM FoodDelivery.type")
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		types := make([]Type, 0)
		for rows.Next() {
			curType := Type {}
			err := rows.Scan(&curType.Id, &curType.Name)
			if err != nil {
				c.String(404, fmt.Sprintf("%v", err))
				return
			}
			types = append(types, curType)
		}
		c.JSON(200, types)
	})

	router.POST("/types", func(c * gin.Context) {
		conn, err := getDBConnection()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		type Body struct {
			NewMealType Type `json:"newMealType"`
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
	})

	router.PUT("/types", func(c * gin.Context) {
		conn, err := getDBConnection()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		type Body struct {
			EditMealType Type `json:"editMealType"`
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
	})

	router.DELETE("/types", func(c * gin.Context) {
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

		tx, err := conn.Begin()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		_, err = tx.Exec("DELETE FROM FoodDelivery.type WHERE id=?", deleteMealTypeId)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}
		_, err = tx.Exec("DELETE FROM FoodDelivery.meals WHERE typeid=?", deleteMealTypeId)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		err = tx.Commit()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		c.String(200, "OK")
	})

	router.GET("/restaurants", func(c *gin.Context) {
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

		restaurants := make([]Restaurant, 0)
		for rows.Next() {
			curRestaurant := Restaurant{}
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
	})

	router.POST("/restaurants", func(c * gin.Context) {
		conn, err := getDBConnection()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		type Body struct {
			NewRestaurant Restaurant `json:"newRestaurant"`
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
	})

	router.PUT("/restaurants", func(c * gin.Context) {
		conn, err := getDBConnection()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		type Body struct {
			EditRestaurant Restaurant `json:"editRestaurant"`
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
	})

	router.DELETE("/restaurants", func(c * gin.Context) {
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
	})

	router.GET("/meals", func(c *gin.Context) {
		conn, err := getDBConnection()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		tx, err := conn.Begin()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		rows, err := tx.Query("SELECT * FROM FoodDelivery.meals")
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		meals := make([]Meal, 0)
		for rows.Next() {
			curMeal := Meal {}
			err := rows.Scan(&curMeal.Id, &curMeal.Name, &curMeal.Description, &curMeal.PhotoURL, &curMeal.Price, &curMeal.TypeId)
			if err != nil {
				c.String(404, fmt.Sprintf("%v", err))
				return
			}

			meals = append(meals, curMeal)
		}

		mealRequests := make([]MealRequest, 0)
		for _, curMeal :=  range meals {
			typeRows, err := tx.Query("SELECT * FROM FoodDelivery.type WHERE id = ? LIMIT 1", curMeal.TypeId)
			if err != nil {
				c.String(404, fmt.Sprintf("%v", err))
				return
			}
			if !typeRows.Next() {
				c.String(404, fmt.Sprintf("Failed to find meal type with id %d", curMeal.TypeId))
				return
			}
			curType := Type{}
			err = typeRows.Scan(&curType.Id, &curType.Name)
			if err != nil {
				c.String(404, fmt.Sprintf("%v", err))
				return
			}

			curMealRequest := MealRequest{
				Id:curMeal.Id,
				Name :curMeal.Name,
				Description:curMeal.Description,
				PhotoURL: curMeal.PhotoURL,
				Price:curMeal.Price,
				Type: curType,

			}

			mealRequests = append(mealRequests, curMealRequest)
		}

		err = tx.Commit()
		if err != nil {
			c.String(404, fmt.Sprintf("AAAAAAA%v", err))
			return
		}

		c.JSON(200, mealRequests)
	})

	router.POST("/meals", func(c * gin.Context) {
		conn, err := getDBConnection()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		type Body struct {
			NewMeal MealRequest `json:"newMeal"`
		}
		b := Body{}
		err = c.BindJSON(&b)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		tx, err := conn.Begin()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		rows, err := tx.Query("SELECT * FROM FoodDelivery.type WHERE id = ?", b.NewMeal.Type.Id)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}
		if !rows.Next() {
			c.String(404, fmt.Sprintf("meal type with id %d doesn't exist", b.NewMeal.Type.Id))
			return
		}

		_, err = tx.Exec("INSERT INTO FoodDelivery.meals(name, description, photourl, price, typeid) VALUES (?)",
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

		err = tx.Commit()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		c.String(200, "OK")
	})

	router.PUT("/meals", func(c * gin.Context) {
		conn, err := getDBConnection()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		type Body struct {
			EditMeal MealRequest `json:"editMeal"`
		}
		b := Body{}
		err = c.BindJSON(&b)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		tx, err := conn.Begin()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		rows, err := tx.Query("SELECT * FROM FoodDelivery.type WHERE id = ?", b.EditMeal.Type.Id)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}
		if !rows.Next() {
			c.String(404, fmt.Sprintf("meal type with id %d doesn't exist", b.EditMeal.Type.Id))
			return
		}

		_, err = tx.Exec("UPDATE FoodDelivery.meals SET name=?, description=?, photourl=?, price=?, typrid=? WHERE id=?",
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

		err = tx.Commit()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		c.String(200, "OK")
	})

	router.DELETE("/meals", func(c * gin.Context) {
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
	})

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
