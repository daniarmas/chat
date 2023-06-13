// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Data interface {
	IsData()
	GetStatus() int
}

type Response interface {
	IsResponse()
	GetStatus() int
	GetMessage() string
	GetError() *Error
	GetData() Data
}

type Chat struct {
	ID           string    `json:"id"`
	Channel      string    `json:"channel"`
	FirstUserID  string    `json:"firstUserId"`
	SecondUserID string    `json:"secondUserId"`
	CreateTime   time.Time `json:"createTime"`
}

type Error struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Details []*ErrorDetails `json:"details,omitempty"`
}

type ErrorDetails struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FetchAllMessagesData struct {
	Status           int        `json:"status"`
	CreateTimeCursor time.Time  `json:"createTimeCursor"`
	Messages         []*Message `json:"messages,omitempty"`
}

func (FetchAllMessagesData) IsData()             {}
func (this FetchAllMessagesData) GetStatus() int { return this.Status }

type FetchAllMessagesInput struct {
	ChatUserID       string     `json:"chatUserId"`
	CreateTimeCursor *time.Time `json:"createTimeCursor,omitempty"`
}

type FetchMessagesResponse struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Error   *Error                `json:"error,omitempty"`
	Data    *FetchAllMessagesData `json:"data,omitempty"`
}

func (FetchMessagesResponse) IsResponse()             {}
func (this FetchMessagesResponse) GetStatus() int     { return this.Status }
func (this FetchMessagesResponse) GetMessage() string { return this.Message }
func (this FetchMessagesResponse) GetError() *Error   { return this.Error }
func (this FetchMessagesResponse) GetData() Data      { return *this.Data }

type GetOrCreateChatData struct {
	Status int   `json:"status"`
	Chat   *Chat `json:"chat"`
}

func (GetOrCreateChatData) IsData()             {}
func (this GetOrCreateChatData) GetStatus() int { return this.Status }

type GetOrCreateChatInput struct {
	ReceiverID string `json:"receiverId"`
}

type GetOrCreateChatResponse struct {
	Status  int                  `json:"status"`
	Message string               `json:"message"`
	Error   *Error               `json:"error,omitempty"`
	Data    *GetOrCreateChatData `json:"data,omitempty"`
}

func (GetOrCreateChatResponse) IsResponse()             {}
func (this GetOrCreateChatResponse) GetStatus() int     { return this.Status }
func (this GetOrCreateChatResponse) GetMessage() string { return this.Message }
func (this GetOrCreateChatResponse) GetError() *Error   { return this.Error }
func (this GetOrCreateChatResponse) GetData() Data      { return *this.Data }

type MeData struct {
	Status int   `json:"status"`
	User   *User `json:"user"`
}

func (MeData) IsData()             {}
func (this MeData) GetStatus() int { return this.Status }

type MeResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Error   *Error  `json:"error,omitempty"`
	Data    *MeData `json:"data,omitempty"`
}

func (MeResponse) IsResponse()             {}
func (this MeResponse) GetStatus() int     { return this.Status }
func (this MeResponse) GetMessage() string { return this.Message }
func (this MeResponse) GetError() *Error   { return this.Error }
func (this MeResponse) GetData() Data      { return *this.Data }

type Message struct {
	ID         string    `json:"id"`
	Content    string    `json:"content"`
	ChatID     string    `json:"chatId"`
	SenderID   string    `json:"senderId"`
	ReceiverID string    `json:"receiverId"`
	CreateTime time.Time `json:"createTime"`
}

type SendMessageData struct {
	Status  int      `json:"status"`
	Message *Message `json:"message"`
}

func (SendMessageData) IsData()             {}
func (this SendMessageData) GetStatus() int { return this.Status }

type SendMessageInput struct {
	Content string `json:"content"`
	ChatID  string `json:"chatId"`
}

type SendMessageResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Error   *Error           `json:"error,omitempty"`
	Data    *SendMessageData `json:"data,omitempty"`
}

func (SendMessageResponse) IsResponse()             {}
func (this SendMessageResponse) GetStatus() int     { return this.Status }
func (this SendMessageResponse) GetMessage() string { return this.Message }
func (this SendMessageResponse) GetError() *Error   { return this.Error }
func (this SendMessageResponse) GetData() Data      { return *this.Data }

type SignInData struct {
	Status       int    `json:"status"`
	User         *User  `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (SignInData) IsData()             {}
func (this SignInData) GetStatus() int { return this.Status }

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Logout   bool   `json:"logout"`
}

type SignInResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   *Error      `json:"error,omitempty"`
	Data    *SignInData `json:"data,omitempty"`
}

func (SignInResponse) IsResponse()             {}
func (this SignInResponse) GetStatus() int     { return this.Status }
func (this SignInResponse) GetMessage() string { return this.Message }
func (this SignInResponse) GetError() *Error   { return this.Error }
func (this SignInResponse) GetData() Data      { return *this.Data }

type SignOutData struct {
	Status int `json:"status"`
}

func (SignOutData) IsData()             {}
func (this SignOutData) GetStatus() int { return this.Status }

type SignOutResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Error   *Error       `json:"error,omitempty"`
	Data    *SignOutData `json:"data,omitempty"`
}

func (SignOutResponse) IsResponse()             {}
func (this SignOutResponse) GetStatus() int     { return this.Status }
func (this SignOutResponse) GetMessage() string { return this.Message }
func (this SignOutResponse) GetError() *Error   { return this.Error }
func (this SignOutResponse) GetData() Data      { return *this.Data }

type Todo struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
	User *User  `json:"user"`
}

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Password   *string   `json:"password,omitempty"`
	Fullname   string    `json:"fullname"`
	Username   string    `json:"username"`
	CreateTime time.Time `json:"createTime"`
}
