package repository

import (
	"context"
	"fmt"
	"transmitter-artemis/config"
	"transmitter-artemis/entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type OutboundRepository interface {
	// SaveSuccess(ctx context.Context, outboundMsg dto.ResponseSuccessQueue, originalRequest []byte) error
	// SaveFailed(ctx context.Context, outboundMsg dto.ResponseFailedQueue, originalRequest []byte) error
	Save(ctx context.Context, clientData entity.ClientData, outboundMessage entity.OutboundMessage) error
}

type outboundRepository struct {
	db *mongo.Client
}

func NewOutboundRepository(db *mongo.Client) *outboundRepository {
	return &outboundRepository{
		db: db,
	}
}

func (or *outboundRepository) Save(ctx context.Context, clientData entity.ClientData, outboundMessage entity.OutboundMessage) error {

	collName := fmt.Sprintf("%s-outbound-msg", clientData.ClientName)

	collection := or.db.Database(config.Configuration.MongoDB.Database).Collection(collName)
	_, err := collection.InsertOne(ctx, outboundMessage)
	if err != nil {
		return err
	}

	return nil
}
