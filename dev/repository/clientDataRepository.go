package repository

import (
	"context"
	"log"
	"transmitter-artemis/config"
	"transmitter-artemis/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClientRepository interface {
	GetAllClientData() ([]entity.ClientData, error)
}

type clientRepository struct {
	Collection *mongo.Collection
}

func NewClientRepository(client *mongo.Client) *clientRepository {
	return &clientRepository{
		Collection: client.Database(config.Configuration.MongoDB.Database).Collection("client-info"),
	}
}

func (cr *clientRepository) GetAllClientData() ([]entity.ClientData, error) {
	var clients []entity.ClientData

	cursor, err := cr.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result entity.ClientData
		if err := cursor.Decode(&result); err != nil {
			log.Println(err)
			continue
		}
		clients = append(clients, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}
