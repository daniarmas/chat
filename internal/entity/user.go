package entity

import (
	"time"
)

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Fullname   string    `json:"fullname"`
	Username   string    `json:"username"`
	CreateTime time.Time `json:"create_time"`
}
