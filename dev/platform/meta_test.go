package platform

import (
	"context"
	"testing"
	"transmitter-artemis/dto"

	"github.com/stretchr/testify/assert"
)

func TestSendRequestToMeta(t *testing.T) {

	t.Run("Test Invalid Token", func(t *testing.T) {
		ctx := context.Background()
		URL := "https://graph.facebook.com"
		token := "token"
		payload := dto.RequestToMeta{

			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               "6283872750005",
			Type:             "text",
			Text: struct {
				PreviewURL bool   `json:"preview_url" bson:"preview_url"`
				Body       string `json:"body" bson:"body"`
			}{
				PreviewURL: false,
				Body:       "Hello, this is a test message",
			},
		}

		app := NewMetaClient()

		_, code, err := app.SendRequestToMeta(ctx, URL, token, payload)

		assert.NotEmpty(t, code)

		assert.NoError(t, err)
	})

}
