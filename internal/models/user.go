package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserOrm struct {
	ID         string    `gorm:"type:uuid;default:uuid_generate_v4()" json:"id" redis:"id"`
	Email      string    `gorm:"unique;not null" json:"email" redis:"email"`
	Password   string    `gorm:"not null" json:"password" redis:"password"`
	Fullname   string    `gorm:"not null" json:"fullname" redis:"fullname"`
	Username   string    `gorm:"unique;not null" json:"username" redis:"username"`
	CreateTime time.Time `json:"create_time" redis:"create_time"`
}

func (UserOrm) TableName() string {
	return "user"
}

func (i *UserOrm) BeforeCreate(tx *gorm.DB) (err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(i.Password), 14)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	i.Password = string(bytes)
	i.CreateTime = time.Now().UTC()
	return
}

// This methods map to and from a UserGorm for avoid using gorm models in the usecases.
func (a *UserOrm) MapToUserGorm(user *entity.User) {
	if user != nil {
		a.ID = user.ID.String()
		a.Email = user.Email
		a.Password = user.Password
		a.Fullname = user.Fullname
		a.Username = user.Username
		a.CreateTime = user.CreateTime
	}
}

func (a UserOrm) MapFromUserGorm() *entity.User {
	var user entity.User
	if a.ID != "" {
		id := uuid.MustParse(a.ID)
		return &entity.User{
			ID:         &id,
			Email:      a.Email,
			Password:   a.Password,
			Fullname:   a.Fullname,
			Username:   a.Username,
			CreateTime: a.CreateTime,
		}
	}
	return &user
}
