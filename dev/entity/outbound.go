package entity

import (
	"time"
	"transmitter-artemis/dto"
)

type OutboundMessage struct {
	MessageID        string               `json:"message_id" bson:"message_id"`
	WAID             string               `json:"wa_id,omitempty" bson:"wa_id,omitempty"`
	To               string               `json:"to" bson:"to"`
	OriginalRequest  string               `json:"original_request" bson:"original_request"`
	Request          dto.RequestToMeta    `json:"request"`
	OriginalResponse dto.ResponseFromMeta `json:"original_response" bson:"original_response"`
	DeliveryReport   []string             `json:"delivery_report" bson:"delivery_report"`
	CreatedAt        time.Time            `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at" bson:"updated_at"`
}
