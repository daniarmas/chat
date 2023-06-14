package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.32

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/daniarmas/chat/graph/model"
	"github.com/daniarmas/chat/internal/inputs"
	"github.com/daniarmas/chat/middleware"
	"github.com/google/uuid"
)

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, input model.SignInInput) (*model.SignInResponse, error) {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	var validationErr = false
	var res model.SignInResponse
	var errorDetails []*model.ErrorDetails

	if input.Email == "" {
		errorDetails = append(errorDetails, &model.ErrorDetails{
			Field:   "email",
			Message: "This field is required",
		})
		validationErr = true
	}

	if input.Email != "" && !emailRegex.MatchString(input.Email) {
		errorDetails = append(errorDetails, &model.ErrorDetails{
			Field:   "email",
			Message: "This field is invalid",
		})
		validationErr = true
	}

	if input.Password == "" {
		errorDetails = append(errorDetails, &model.ErrorDetails{
			Field:   "password",
			Message: "This field is required",
		})
		validationErr = true
	}

	if validationErr {
		res.Message = "Bad request"
		res.Status = http.StatusBadRequest
		res.Data = nil
		res.Error = &model.Error{
			Code:    "INVALID_ARGUMENT",
			Message: "The request contains invalid arguments",
			Details: errorDetails,
		}
		return &res, nil
	}

	result, err := r.AuthUsecase.SignIn(ctx, inputs.SignInInput{Email: input.Email, Password: input.Password, Logout: input.Logout})
	if err != nil {
		switch err.Error() {
		case "the credentials are incorrect":
			res.Message = http.StatusText(http.StatusUnauthorized)
			res.Status = http.StatusUnauthorized
			res.Data = nil
			res.Error = &model.Error{
				Code:    "INVALID_CREDENTIALS",
				Message: "The credentials are incorrect.",
				Details: nil,
			}
			return &res, nil
		case "the user is already logged in":
			res.Message = http.StatusText(http.StatusConflict)
			res.Status = http.StatusConflict
			res.Data = nil
			res.Error = &model.Error{
				Code:    "USER_ALREADY_LOGGED_IN",
				Message: "The user is already logged in.",
				Details: nil,
			}
			return &res, nil
		default:
			res.Message = http.StatusText(http.StatusInternalServerError)
			res.Status = http.StatusInternalServerError
			res.Data = nil
			res.Error = &model.Error{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "The server has an internal error.",
				Details: nil,
			}
			return &res, nil
		}
	}

	res.Message = "Success"
	res.Status = http.StatusOK
	res.Data = &model.SignInData{
		Status: http.StatusOK,
		User: &model.User{
			ID:         result.User.ID.String(),
			Email:      result.User.Email,
			Fullname:   result.User.Fullname,
			Username:   result.User.Username,
			CreateTime: result.User.CreateTime,
		},
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	return &res, nil
}

// SignOut is the resolver for the signOut field.
func (r *mutationResolver) SignOut(ctx context.Context) (*model.SignOutResponse, error) {
	var res model.SignOutResponse

	user := middleware.ForContext(ctx)
	if user == nil {
		res.Message = http.StatusText(http.StatusUnauthorized)
		res.Status = http.StatusUnauthorized
		res.Data = nil
		res.Error = &model.Error{
			Code:    "ACCESS_TOKEN_MISSING",
			Message: "This request requires an access token. Please provide a valid access token and try again.",
			Details: nil,
		}
		return &res, nil
	}

	err := r.AuthUsecase.SignOut(ctx, user.ID.String())
	if err != nil {
		switch err.Error() {
		default:
			res.Message = http.StatusText(http.StatusInternalServerError)
			res.Status = http.StatusInternalServerError
			res.Data = nil
			res.Error = &model.Error{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "The server has an internal error.",
				Details: nil,
			}
			return &res, nil
		}
	}

	res.Message = http.StatusText(http.StatusNoContent)
	res.Status = http.StatusNoContent
	res.Data = nil
	res.Error = nil
	return &res, nil
}

// SendMessage is the resolver for the sendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, input model.SendMessageInput) (*model.SendMessageResponse, error) {
	var res model.SendMessageResponse
	var errorDetails []*model.ErrorDetails
	var validationErr = false
	var chatId uuid.UUID

	user := middleware.ForContext(ctx)
	if user == nil {
		res.Message = http.StatusText(http.StatusUnauthorized)
		res.Status = http.StatusUnauthorized
		res.Data = nil
		res.Error = &model.Error{
			Code:    "ACCESS_TOKEN_MISSING",
			Message: "This request requires an access token. Please provide a valid access token and try again.",
			Details: nil,
		}
		return &res, nil
	}

	if input.Content == "" {
		errorDetails = append(errorDetails, &model.ErrorDetails{
			Field:   "content",
			Message: "This field is required",
		})
		validationErr = true
	}

	if input.ChatID == "" {
		errorDetails = append(errorDetails, &model.ErrorDetails{
			Field:   "receiver_id",
			Message: "This field is required",
		})
		validationErr = true
	} else {
		chatId = uuid.MustParse(input.ChatID)
	}

	if validationErr {
		res.Message = "Bad request"
		res.Status = http.StatusBadRequest
		res.Data = nil
		res.Error = &model.Error{
			Code:    "INVALID_ARGUMENT",
			Message: "The request contains invalid arguments",
			Details: errorDetails,
		}
		return &res, nil
	}

	result, err := r.MessageUsecase.SendMessage(ctx, inputs.SendMessage{ChatID: &chatId, Content: input.Content}, user.ID)
	if err != nil {
		switch err.Error() {
		default:
			res.Message = http.StatusText(http.StatusInternalServerError)
			res.Status = http.StatusInternalServerError
			res.Data = nil
			res.Error = &model.Error{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "The server has an internal error.",
				Details: nil,
			}
			return &res, nil
		}
	}

	res.Message = "Success"
	res.Status = http.StatusOK
	res.Data = &model.SendMessageData{
		Message: &model.Message{ID: result.ID.String(), Content: result.Content, ChatID: result.ChatId.String(), CreateTime: result.CreateTime},
	}
	res.Error = nil

	return &res, nil
}

// GetOrCreateChat is the resolver for the getOrCreateChat field.
func (r *mutationResolver) GetOrCreateChat(ctx context.Context, input model.GetOrCreateChatInput) (*model.GetOrCreateChatResponse, error) {
	var res model.GetOrCreateChatResponse
	var errorDetails []*model.ErrorDetails
	var validationErr = false
	var receiverId uuid.UUID

	user := middleware.ForContext(ctx)
	if user == nil {
		res.Message = http.StatusText(http.StatusUnauthorized)
		res.Status = http.StatusUnauthorized
		res.Data = nil
		res.Error = &model.Error{
			Code:    "ACCESS_TOKEN_MISSING",
			Message: "This request requires an access token. Please provide a valid access token and try again.",
			Details: nil,
		}
		return &res, nil
	}

	if input.ReceiverID == "" {
		errorDetails = append(errorDetails, &model.ErrorDetails{
			Field:   "receiver_id",
			Message: "This field is required",
		})
		validationErr = true
	} else {
		receiverId = uuid.MustParse(input.ReceiverID)
	}

	if validationErr {
		res.Message = "Bad request"
		res.Status = http.StatusBadRequest
		res.Data = nil
		res.Error = &model.Error{
			Code:    "INVALID_ARGUMENT",
			Message: "The request contains invalid arguments",
			Details: errorDetails,
		}
		return &res, nil
	}

	result, err := r.ChatUsecase.GetOrCreateChat(ctx, inputs.GetOrCreateChatInput{ReceiverId: &receiverId}, user.ID.String())
	if err != nil {
		switch err.Error() {
		default:
			res.Message = http.StatusText(http.StatusInternalServerError)
			res.Status = http.StatusInternalServerError
			res.Data = nil
			res.Error = &model.Error{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "The server has an internal error.",
				Details: nil,
			}
			return &res, nil
		}
	}

	res.Message = "Success"
	res.Status = http.StatusOK
	res.Data = &model.GetOrCreateChatData{
		Chat: &model.Chat{ID: result.ID.String(), FirstUserID: result.FirstUserId.String(), SecondUserID: result.SecondUserId.String(), CreateTime: result.CreateTime},
	}
	res.Error = nil

	return &res, nil
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.MeResponse, error) {
	var res model.MeResponse

	user := middleware.ForContext(ctx)
	if user == nil {
		res.Message = http.StatusText(http.StatusUnauthorized)
		res.Status = http.StatusUnauthorized
		res.Data = nil
		res.Error = &model.Error{
			Code:    "ACCESS_TOKEN_MISSING",
			Message: "This request requires an access token. Please provide a valid access token and try again.",
			Details: nil,
		}
		return &res, nil
	}

	result, err := r.AuthUsecase.Me(ctx, user.ID.String())
	if err != nil {
		switch err.Error() {
		default:
			res.Message = http.StatusText(http.StatusInternalServerError)
			res.Status = http.StatusInternalServerError
			res.Data = nil
			res.Error = &model.Error{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "The server has an internal error.",
				Details: nil,
			}
			return &res, nil
		}
	}

	res.Message = "Success"
	res.Status = http.StatusOK
	res.Data = &model.MeData{
		User: &model.User{ID: result.ID.String(), Email: result.Email, Fullname: result.Fullname, Username: result.Username, CreateTime: result.CreateTime},
	}
	res.Error = nil

	return &res, nil
}

// FetchMessages is the resolver for the fetchMessages field.
func (r *queryResolver) FetchMessages(ctx context.Context, input model.FetchAllMessagesInput) (*model.FetchMessagesResponse, error) {
	var res model.FetchMessagesResponse
	messages := make([]*model.Message, 0)
	var createTimeCursor time.Time

	if input.CreateTimeCursor != nil && !input.CreateTimeCursor.IsZero() {
		createTimeCursor = input.CreateTimeCursor.UTC()
	}

	user := middleware.ForContext(ctx)
	if user == nil {
		res.Message = http.StatusText(http.StatusUnauthorized)
		res.Status = http.StatusUnauthorized
		res.Data = nil
		res.Error = &model.Error{
			Code:    "ACCESS_TOKEN_MISSING",
			Message: "This request requires an access token. Please provide a valid access token and try again.",
			Details: nil,
		}
		return &res, nil
	}

	result, err := r.MessageUsecase.GetMessageByChat(ctx, inputs.GetMessagesByChatId{ChatId: input.ChatID}, user.ID.String(), createTimeCursor)
	if err != nil {
		switch err.Error() {
		default:
			res.Message = http.StatusText(http.StatusInternalServerError)
			res.Status = http.StatusInternalServerError
			res.Data = nil
			res.Error = &model.Error{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "The server has an internal error.",
				Details: nil,
			}
			return &res, nil
		}
	}

	for _, element := range result.Messages {
		messages = append(messages, &model.Message{
			ID:         element.ID.String(),
			Content:    element.Content,
			ChatID:     element.ChatId.String(),
			CreateTime: element.CreateTime,
		})
	}

	if len(messages) == 11 {
		messages = messages[:len(messages)-1]
	}

	var createTimeCursorRes time.Time

	if len(messages) != 0 {
		createTimeCursorRes = messages[len(messages)-1].CreateTime
	}

	res.Message = "Success"
	res.Status = http.StatusOK
	res.Data = &model.FetchAllMessagesData{
		Messages:         messages,
		CreateTimeCursor: &createTimeCursorRes,
	}
	res.Error = nil

	return &res, nil
}

// FetchChats is the resolver for the fetchChats field.
func (r *queryResolver) FetchChats(ctx context.Context, input model.FetchAllChatsInput) (*model.FetchChatsResponse, error) {
	var res model.FetchChatsResponse
	chats := make([]*model.Chat, 0)
	var updateTimeCursor time.Time

	if input.UpdateTimeCursor != nil && !input.UpdateTimeCursor.IsZero() {
		updateTimeCursor = input.UpdateTimeCursor.UTC()
	}

	user := middleware.ForContext(ctx)
	if user == nil {
		res.Message = http.StatusText(http.StatusUnauthorized)
		res.Status = http.StatusUnauthorized
		res.Data = nil
		res.Error = &model.Error{
			Code:    "ACCESS_TOKEN_MISSING",
			Message: "This request requires an access token. Please provide a valid access token and try again.",
			Details: nil,
		}
		return &res, nil
	}

	result, err := r.ChatUsecase.GetChats(ctx, user.ID.String(), updateTimeCursor)
	if err != nil {
		switch err.Error() {
		default:
			res.Message = http.StatusText(http.StatusInternalServerError)
			res.Status = http.StatusInternalServerError
			res.Data = nil
			res.Error = &model.Error{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "The server has an internal error.",
				Details: nil,
			}
			return &res, nil
		}
	}

	for _, element := range result.Chats {
		chats = append(chats, &model.Chat{
			ID:           element.ID.String(),
			FirstUserID:  element.FirstUserId.String(),
			SecondUserID: element.SecondUserId.String(),
			CreateTime:   element.CreateTime,
		})
	}

	if len(chats) == 11 {
		chats = chats[:len(chats)-1]
	}

	var updateTimeCursorRes time.Time

	if len(chats) != 0 {
		updateTimeCursorRes = chats[len(chats)-1].CreateTime
	}

	res.Message = "Success"
	res.Status = http.StatusOK
	res.Data = &model.FetchChatsData{
		Chats:            chats,
		UpdateTimeCursor: &updateTimeCursorRes,
	}
	res.Error = nil

	return &res, nil
}

// CurrentTime is the resolver for the currentTime field.
func (r *subscriptionResolver) CurrentTime(ctx context.Context) (<-chan *time.Time, error) {
	// First you'll need to `make()` your channel. Use your type here!
	ch := make(chan *time.Time)

	// You can (and probably should) handle your channels in a central place outside of `schema.resolvers.go`.
	// For this example we'll simply use a Goroutine with a simple loop.
	go func() {
		for {
			// In our example we'll send the current time every second.
			time.Sleep(1 * time.Second)
			fmt.Println("Tick")

			// Prepare your object.
			currentTime := time.Now()

			// The channel may have gotten closed due to the client disconnecting.
			// To not have our Goroutine block or panic, we do the send in a select block.
			// This will jump to the default case if the channel is closed.
			select {
			case ch <- &currentTime: // This is the actual send.
				// Our message went through, do nothing
			default: // This is run when our send does not work.
				fmt.Println("Channel closed.")
				// You can handle any deregistration of the channel here.
				return // We'll just return ending the routine.
			}
		}
	}()

	// We return the channel and no error.
	return ch, nil
}

// ReceiveMessages is the resolver for the receiveMessages field.
func (r *subscriptionResolver) ReceiveMessages(ctx context.Context, input model.ReceiveMessagesInput) (<-chan *model.Message, error) {
	res := make(chan *model.Message)

	user := middleware.ForContext(ctx)
	if user == nil {
		return nil, errors.New("access token missing")
	}

	result, err := r.MessageUsecase.ReceiveMessages(ctx, inputs.ReceiveMessagesInput{ChatId: input.ChatID})
	if err != nil {
		return nil, errors.New("internal server error")
	}

	// goroutine for publishing model.Message objects to the publishing channel
	go func() {
		defer close(res)

		for entityMsg := range result {
			// convert the entity.Message to model.Message
			modelMsg := &model.Message{
				ID:         entityMsg.ID.String(),
				Content:    entityMsg.Content,
				ChatID:     entityMsg.ID.String(),
				UserID:     entityMsg.UserId.String(),
				CreateTime: entityMsg.CreateTime,
			}

			// send the model.Message to the modelMessages channel
			if modelMsg.UserID != user.ID.String() {
				res <- modelMsg
			}
		}
	}()

	return res, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}
