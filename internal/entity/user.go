package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         *uuid.UUID `json:"id"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Fullname   string     `json:"fullname"`
	Username   string     `json:"username"`
	CreateTime time.Time  `json:"create_time"`
}
