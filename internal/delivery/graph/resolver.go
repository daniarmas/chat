package graph

import "github.com/daniarmas/chat/internal/usecases"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AuthUsecase    usecases.AuthUsecase
	MessageUsecase usecases.MessageUsecase
	ChatUsecase    usecases.ChatUsecase
}
