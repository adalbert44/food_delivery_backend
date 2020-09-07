package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Type struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
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
		_, err := getDBConnection()
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		input := make([]byte, 1024)
		num, err := c.Request.Body.Read(input)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}
		c.String(200, string(input[:num]))

		/*
		type Body struct {
			DeleteMealTypeId int `json:"deleteMealTypeId"`
		}
		b := Body{}
		err = c.BindJSON(&b)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		_, err = conn.Exec("DELETE FROM FoodDelivery.type WHERE id=?", b.DeleteMealTypeId)
		if err != nil {
			c.String(404, fmt.Sprintf("%v", err))
			return
		}

		c.String(200, "OK")
		*/
	})

	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
