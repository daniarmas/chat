package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.32

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	"github.com/daniarmas/chat/graph/model"
	"github.com/daniarmas/chat/internal/inputs"
)

// SignIn is the resolver for the signIn field.
func (r *mutationResolver) SignIn(ctx context.Context, in model.SignInInput) (*model.SignInResponse, error) {
	var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	var validationErr = false
	var res model.SignInResponse
	var errorDetails []*model.ErrorDetails

	if in.Email == "" {
		errorDetails = append(errorDetails, &model.ErrorDetails{
			Field:   "email",
			Message: "This field is required",
		})
		validationErr = true
	}

	if in.Email != "" && !emailRegex.MatchString(in.Email) {
		errorDetails = append(errorDetails, &model.ErrorDetails{
			Field:   "email",
			Message: "This field is invalid",
		})
		validationErr = true
	}

	if in.Password == "" {
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

	result, err := r.AuthUsecase.SignIn(ctx, inputs.SignInInput{Email: in.Email, Password: in.Password, Logout: in.Logout})
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
			Password:   "",
			Fullname:   result.User.Fullname,
			Username:   result.User.Username,
			CreateTime: result.User.CreateTime,
		},
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}

	// c.JSON(http.StatusOK, res)
	return &res, nil
}

// Todos is the resolver for the todos field.
func (r *queryResolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
