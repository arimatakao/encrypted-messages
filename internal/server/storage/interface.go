package storage

import "context"

type UserStorager interface {
	AddUser(u *User) error
	ReadUser(username string) (*User, error)
	DeleteUser(username string) error
}

type MessageStorager interface {
	AddMessage(m *MessageReq) (string, error)
	ReadMessage(id string) (*Message, error)
	ReadAllMessages() ([]*Message, error)
	DeleteMessage(id string) error
	Disconnect(ctx context.Context) error
}

type Storager interface {
	MessageStorager
	UserStorager
}

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	Id       string `json:"id" bson:"_id,omitempty"`
	OwnerId  string `json:"owner_id" bson:"owner_id"`
	IsPublic bool   `json:"is_public" bson:"is_public"`
	Password string `json:"password" bson:"password"`
	Content  string `json:"content" bsot:"content"`
}

type MessageReq struct {
	OwnerId  *string `json:"owner_id" bson:"owner_id"`
	IsPublic *bool   `json:"is_public" bson:"is_public"`
	Content  *string `json:"content" bsot:"content"`
}

func (m MessageReq) SomeEmpty() bool {
	return m.Content == nil ||
		m.IsPublic == nil
}
