package response

import (
	"github.com/daniarmas/chat/internal/entity"
)

type SignInBusinessResponse struct {
	User         *entity.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}
