package models

import (
	"time"

	"github.com/daniarmas/chat/internal/entity"
	bcryptutils "github.com/daniarmas/chat/pkg/bcrypt_utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UserOrm struct {
	ID         *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Email      string     `gorm:"unique;not null" json:"email"`
	Password   string     `gorm:"not null" json:"password"`
	Fullname   string     `gorm:"not null" json:"fullname"`
	Username   string     `gorm:"unique;not null" json:"username"`
	CreateTime time.Time  `json:"create_time"`
}

func (UserOrm) TableName() string {
	return "user"
}

func (i *UserOrm) BeforeCreate(tx *gorm.DB) (err error) {
	hash, err := bcryptutils.HashPassword(i.Password)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	i.Password = hash
	i.CreateTime = time.Now().UTC()
	return
}

// This methods map to and from a UserGorm for avoid using gorm models in the usecases.
func (a *UserOrm) MapToUserGorm(user *entity.User) {
	if user != nil {
		a.ID = user.ID
		a.Email = user.Email
		a.Password = user.Password
		a.Fullname = user.Fullname
		a.Username = user.Username
		a.CreateTime = user.CreateTime
	}
}

func (a UserOrm) MapFromUserGorm() *entity.User {
	return &entity.User{
		ID:         a.ID,
		Email:      a.Email,
		Password:   a.Password,
		Fullname:   a.Fullname,
		Username:   a.Username,
		CreateTime: a.CreateTime,
	}
}
