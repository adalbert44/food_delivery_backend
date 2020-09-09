package models

type Meal struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	PhotoURL string `json:"photoUrl"`
	Price int `json:"price"`
	TypeId int `json:"type"`
}

type MealRequest struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	PhotoURL string `json:"photoUrl"`
	Price int `json:"price"`
	Type Type `json:"type"`
}

