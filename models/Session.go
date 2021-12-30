package models

import "time"

type Session struct {
	Id             string    `json:"id" gorm:"primary_key"`
	UserId         int       `json:"user_id" gorm:"not null;unique"`
	ExpirationTime time.Time `json:"expiration_time" gorm:"not null"`
}
