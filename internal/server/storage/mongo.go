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
	client *mongo.Client
	coll   *mongo.Collection
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

	coll := client.Database("encmsg").Collection("messages")

	return &mongoDB{
		client: client,
		coll:   coll,
	}, nil
}

func (m *mongoDB) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

func (m mongoDB) AddMessage(msg *MessageReq) (string, error) {
	res, err := m.coll.InsertOne(context.TODO(), msg)
	if err != nil {
		return "", err
	}
	resId, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("failed convert _id to string")
	}

	return resId.Hex(), nil
}

func (m mongoDB) ReadMessage(id string) (*Message, error) {
	obj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("failed convert string to _id: %s", id)
	}

	filter := bson.D{{Key: "_id", Value: obj}}

	res := m.coll.FindOne(context.TODO(), filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var msg Message
	if err := res.Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil

}

func (m mongoDB) ReadAllMessages() ([]*Message, error) {
	return nil, nil
}

func (m mongoDB) DeleteMessage(id string) error {
	return nil
}
