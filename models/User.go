package models

type User struct {
	Id       int    `json:"id" gorm:"primary_key"`
	Email    string `json:"email" gorm:"not null;unique"`
	Password string `json:"password"`
	Status   string `json:"status"`
}
