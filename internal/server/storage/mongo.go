package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoDB struct {
	client    *mongo.Client
	usersColl *mongo.Collection
	msgColl   *mongo.Collection
}

func NewMongoDB(connUrl string) (*mongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connUrl))
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	usersColl := client.Database("encmsg").Collection("users")
	msgColl := client.Database("encmsg").Collection("messages")

	return &mongoDB{
		client:    client,
		usersColl: usersColl,
		msgColl:   msgColl,
	}, nil
}

func (m *mongoDB) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m mongoDB) AddUser(u *UserReq) error {
	_, err := m.usersColl.InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}

	return nil
}

func (m mongoDB) ReadUser(u *UserReq) (User, error) {
	res := m.usersColl.FindOne(context.TODO(), u)
	if res.Err() != nil {
		return User{}, res.Err()
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return User{}, err
	}

	return user, nil
}

func (m mongoDB) ReadUserByUsername(username string) (*User, error) {
	filter := bson.D{{Key: "username", Value: username}}

	res := m.usersColl.FindOne(context.TODO(), filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m mongoDB) ReadUserById(id string) (*User, error) {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to _id: %s", id)
	}

	filter := bson.D{{Key: "_id", Value: obj}}

	res := m.usersColl.FindOne(context.TODO(), filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var user User
	if err := res.Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m mongoDB) DeleteUser(id string) error {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert string to _id: %s", id)
	}

	filter := bson.D{{Key: "_id", Value: obj}}

	_, err = m.usersColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete message with _id: %s", id)
	}

	return nil
}

func (m mongoDB) AddMessage(msgReq *Message) (string, error) {
	msg := MessageReq{
		IsPublic: &msgReq.IsPublic,
		Content:  &msgReq.Content,
	}
	res, err := m.msgColl.InsertOne(context.TODO(), msg)
	if err != nil {
		return "", err
	}
	resId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed to convert _id to string")
	}

	return resId.Hex(), nil
}

func (m mongoDB) ReadMessage(id string) (*Message, error) {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to _id: %s", id)
	}

	filter := bson.D{{Key: "_id", Value: obj}}

	res := m.msgColl.FindOne(context.TODO(), filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var msg Message
	if err := res.Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil

}

func (m mongoDB) ReadAllMessages(owner_id string) ([]*Message, error) {
	return nil, nil
}

func (m mongoDB) DeleteMessage(id string) error {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert string to _id: %s", id)
	}

	filter := bson.D{{Key: "_id", Value: obj}}

	_, err = m.msgColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete message with _id: %s", id)
	}

	return nil
}

func (m mongoDB) DeleteAllMessages(owner_id string) error {
	return nil
}
