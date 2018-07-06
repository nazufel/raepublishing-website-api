package models

//Struct to hold user data
type User struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
	Id     string `json:"id"`
}
