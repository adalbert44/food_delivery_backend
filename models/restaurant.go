package models

type Restaurant struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Location string `json:"location"`
	PhotoURL string `json:"photoUrl"`
}

