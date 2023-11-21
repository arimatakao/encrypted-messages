package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStorager interface {
	AddUser(u *UserReq) error
	ReadUser(u *UserReq) (User, error)
	ReadUserByUsername(username string) (*User, error)
	ReadUserById(id string) (*User, error)
	DeleteUser(id string) error
}

type MessageStorager interface {
	AddMessage(m *Message) (string, error)
	ReadMessage(id string) (*Message, error)
	ReadAllMessages(owner_id string) ([]*Message, error)
	DeleteMessage(id string) error
	DeleteAllMessages(owner_id string) error
	Disconnect(ctx context.Context) error
}

type Storager interface {
	MessageStorager
	UserStorager
}

type User struct {
	Id       string `json:"id" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type UserReq struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Message struct {
	Id       string `json:"id" bson:"_id"`
	OwnerId  string `json:"owner_id" bson:"owner_id"`
	IsPublic bool   `json:"is_public" bson:"is_public"`
	Content  string `json:"content" bson:"content"`
	Password string `json:"password" bson:"password"`
}

type MessageReq struct {
	OwnerId  primitive.ObjectID `json:"owner_id" bson:"owner_id"`
	IsPublic *bool              `json:"is_public" bson:"is_public"`
	Content  *string            `json:"content" bson:"content"`
}

func (m Message) IsEmpty() bool {
	return m.Content == ""
}
