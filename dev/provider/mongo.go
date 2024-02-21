package provider

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"transmitter-artemis/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDBClient creates a new MongoDB client
func NewMongoDBClient() (*mongo.Client, error) {
	cfg := config.Configuration.MongoDB

	optionsStr := fmt.Sprintf("&%s=%s", "authSource", url.QueryEscape(cfg.AuthSource))

	// Set the MongoDB client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d/?%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, optionsStr))

	// Connect to the MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return client, nil
}
