package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
)

type User struct {
	ID         string    `json:"id" redis:"id"`
	Email      string    `json:"email" redis:"email"`
	Password   string    `json:"password" redis:"password"`
	Fullname   string    `json:"fullname" redis:"fullname"`
	Username   string    `json:"username" redis:"username"`
	CreateTime time.Time `json:"create_time" redis:"create_time"`
}

// This methods map to and from a UserGorm for avoid using gorm models in the usecases.
func (a *User) MapToUserModel(user *entity.User) {
	if user != nil {
		a.ID = user.ID
		a.Email = user.Email
		a.Password = user.Password
		a.Fullname = user.Fullname
		a.Username = user.Username
		a.CreateTime = user.CreateTime
	}
}

func (a User) MapFromUserModel() *entity.User {
	var user entity.User
	if a.ID != "" {
		return &entity.User{
			ID:         a.ID,
			Email:      a.Email,
			Password:   a.Password,
			Fullname:   a.Fullname,
			Username:   a.Username,
			CreateTime: a.CreateTime,
		}
	}
	return &user
}
