package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
	"transmitter-artemis/dto"
	"transmitter-artemis/entity"
	mocksPlatform "transmitter-artemis/mocks/platform"
	mocksProvider "transmitter-artemis/mocks/provider"
	mocksRepo "transmitter-artemis/mocks/repository"
	"transmitter-artemis/provider"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSendMessage_all(t *testing.T) {
	type want struct {
		err error
	}
	type test struct {
		name       string
		msgString  string
		clientData entity.ClientData
		// requestBodyMeta  dto.RequestToMeta
		// responseBodyMeta dto.ResponseFromMeta

		outboundRepoMock func() (outboundRepoMock *mocksRepo.OutboundRepository)
		drRepoMock       func() (drRepoMock *mocksRepo.DRRepository)
		loggerMock       func() (loggerMock *mocksProvider.ILogger)
		metaPlatformMock func() (metaPlatformMock *mocksPlatform.MetaClient)
		want             want
	}

	ctx := context.Background()
	RequestToMeta := dto.RequestToMeta{
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

	responseFromMeta := dto.ResponseFromMeta{
		MessagingProduct: "whatsapp",
		Contacts: []struct {
			Input string `json:"input" bson:"input"`
			WAID  string `json:"wa_id,omitempty" bson:"wa_id,omitempty"`
		}{
			{Input: "6283872750005", WAID: "6283872750005"},
		},
		Messages: []struct {
			ID string `json:"id"`
		}{
			{ID: "wamid.abc123"},
		},
	}

	clientData := entity.ClientData{
		ClientName:    "lino",
		Token:         "abc",
		PhoneNumberID: "123",
		WAHost:        "https://graph.facebook.com",
	}

	URL := fmt.Sprintf("%s/%s/messages", clientData.WAHost, clientData.PhoneNumberID)

	rightNow := time.Now()
	timeString := fmt.Sprintf("%d", rightNow.Unix())
	drMsg := fmt.Sprintf("message_id=1&wa_id=wamid.abc123&deliverystatus=sent&time=%v", timeString)

	tests := []test{

		{name: "Test Success SendMessage",
			msgString: "message_id=1&to=6283872750005&type=text&text[preview_url]=false&text[body]=Hello, this is a test message",
			clientData: entity.ClientData{
				ClientName:    "lino",
				Token:         "abc",
				PhoneNumberID: "123",
				WAHost:        "https://graph.facebook.com",
			},
			metaPlatformMock: func() (metaPlatformMock *mocksPlatform.MetaClient) {
				metaPlatformMock = mocksPlatform.NewMetaClient(t)

				metaPlatformMock.On("SendRequestToMeta", ctx, URL, clientData.Token, RequestToMeta).Return(responseFromMeta, 200, nil)
				return
			},

			drRepoMock: func() (drRepoMock *mocksRepo.DRRepository) {
				drRepoMock = mocksRepo.NewDRRepository(t)
				drRepoMock.On("Produce", ctx, clientData, drMsg).Return(nil)
				return
			},
			outboundRepoMock: func() (outboundRepoMock *mocksRepo.OutboundRepository) {
				outboundRepoMock = mocksRepo.NewOutboundRepository(t)
				outboundRepoMock.On("Save", ctx, clientData, mock.Anything).Return(nil)
				return
			},
			loggerMock: func() (loggerMock *mocksProvider.ILogger) {
				loggerMock = mocksProvider.NewILogger(t)
				loggerMock.On("Infof", provider.AppLog, mock.Anything)
				return
			},
			want: want{
				err: nil,
			},
		},
		{name: "Test Failed SendMessage",
			msgString: "message_id=1&to=6283872750005&type=text&text[preview_url]=false&text[body]=Hello, this is a test message",
			clientData: entity.ClientData{
				ClientName:    "lino",
				Token:         "abc",
				PhoneNumberID: "123",
				WAHost:        "https://graph.facebook.com",
			},
			metaPlatformMock: func() (metaPlatformMock *mocksPlatform.MetaClient) {
				metaPlatformMock = mocksPlatform.NewMetaClient(t)
				errorRes := dto.ErrorRes{
					Message: "Error message",
				}

				responseFailed := dto.ResponseFromMeta{
					Error: &errorRes,
				}
				metaPlatformMock.On("SendRequestToMeta", ctx, URL, clientData.Token, RequestToMeta).Return(responseFailed, 200, nil)
				return
			},

			drRepoMock: func() (drRepoMock *mocksRepo.DRRepository) {
				drRepoMock = mocksRepo.NewDRRepository(t)
				drRepoMock.On("Produce", ctx, clientData, mock.Anything).Return(nil)
				return
			},
			outboundRepoMock: func() (outboundRepoMock *mocksRepo.OutboundRepository) {
				outboundRepoMock = mocksRepo.NewOutboundRepository(t)
				outboundRepoMock.On("Save", ctx, clientData, mock.Anything).Return(nil)
				return
			},
			loggerMock: func() (loggerMock *mocksProvider.ILogger) {
				loggerMock = mocksProvider.NewILogger(t)
				loggerMock.On("Infof", provider.AppLog, mock.Anything)
				return
			},
			want: want{
				err: nil,
			},
		},
		{
			name:      "Test Invalid Data msg from Queue",
			msgString: "to=6283872750005&type;=text", // <-- INVALID INPUT
			clientData: entity.ClientData{
				ClientName:    "lino",
				Token:         "abc",
				PhoneNumberID: "123",
				WAHost:        "https://graph.facebook.com",
			},
			metaPlatformMock: func() (metaPlatformMock *mocksPlatform.MetaClient) {
				metaPlatformMock = mocksPlatform.NewMetaClient(t)
				return
			},

			drRepoMock: func() (drRepoMock *mocksRepo.DRRepository) {
				drRepoMock = mocksRepo.NewDRRepository(t)
				return
			},
			outboundRepoMock: func() (outboundRepoMock *mocksRepo.OutboundRepository) {
				outboundRepoMock = mocksRepo.NewOutboundRepository(t)
				return
			},
			loggerMock: func() (loggerMock *mocksProvider.ILogger) {
				loggerMock = mocksProvider.NewILogger(t)
				loggerMock.On("Errorf", provider.AppLog, mock.Anything)
				return
			},
			want: want{
				err: errors.New("invalid semicolon separator in query"),
			},
		},
		{
			name: "Test Cant Send request to meta",
			clientData: entity.ClientData{
				ClientName:    "lino",
				Token:         "abc",
				PhoneNumberID: "123",
				WAHost:        "https://graph.facebook.com",
			},
			metaPlatformMock: func() (metaPlatformMock *mocksPlatform.MetaClient) {
				metaPlatformMock = mocksPlatform.NewMetaClient(t)
				metaPlatformMock.On("SendRequestToMeta", ctx, URL, clientData.Token, mock.Anything).Return(dto.ResponseFromMeta{}, 500, errors.New("Cannot Send Request"))
				return
			},

			drRepoMock: func() (drRepoMock *mocksRepo.DRRepository) {
				drRepoMock = mocksRepo.NewDRRepository(t)
				return
			},
			outboundRepoMock: func() (outboundRepoMock *mocksRepo.OutboundRepository) {
				outboundRepoMock = mocksRepo.NewOutboundRepository(t)
				return
			},
			loggerMock: func() (loggerMock *mocksProvider.ILogger) {
				loggerMock = mocksProvider.NewILogger(t)
				loggerMock.On("Errorf", provider.AppLog, mock.Anything)
				return
			},
			want: want{
				err: errors.New("Cannot Send Request"),
			},
		},
		{
			name: "Test Cannot Save DR queue to Artemis",
			clientData: entity.ClientData{
				ClientName:    "lino",
				Token:         "abc",
				PhoneNumberID: "123",
				WAHost:        "https://graph.facebook.com",
			},
			metaPlatformMock: func() (metaPlatformMock *mocksPlatform.MetaClient) {
				metaPlatformMock = mocksPlatform.NewMetaClient(t)
				metaPlatformMock.On("SendRequestToMeta", ctx, URL, clientData.Token, mock.Anything).Return(responseFromMeta, 200, nil)
				return
			},

			drRepoMock: func() (drRepoMock *mocksRepo.DRRepository) {
				drRepoMock = mocksRepo.NewDRRepository(t)
				drRepoMock.On("Produce", ctx, clientData, mock.Anything).Return(errors.New("Cannot Save to DR-queue Artemis"))
				return
			},
			outboundRepoMock: func() (outboundRepoMock *mocksRepo.OutboundRepository) {
				outboundRepoMock = mocksRepo.NewOutboundRepository(t)
				return
			},
			loggerMock: func() (loggerMock *mocksProvider.ILogger) {
				loggerMock = mocksProvider.NewILogger(t)
				loggerMock.On("Errorf", provider.AppLog, mock.Anything)
				return
			},
			want: want{
				err: errors.New("Cannot Save to DR-queue Artemis"),
			},
		},
		{
			name: "Test Cannot Store msg to Mongo",
			clientData: entity.ClientData{
				ClientName:    "lino",
				Token:         "abc",
				PhoneNumberID: "123",
				WAHost:        "https://graph.facebook.com",
			},
			metaPlatformMock: func() (metaPlatformMock *mocksPlatform.MetaClient) {
				metaPlatformMock = mocksPlatform.NewMetaClient(t)
				metaPlatformMock.On("SendRequestToMeta", ctx, URL, clientData.Token, mock.Anything).Return(responseFromMeta, 200, nil)
				return
			},

			drRepoMock: func() (drRepoMock *mocksRepo.DRRepository) {
				drRepoMock = mocksRepo.NewDRRepository(t)
				drRepoMock.On("Produce", ctx, clientData, mock.Anything).Return(nil)
				return
			},
			outboundRepoMock: func() (outboundRepoMock *mocksRepo.OutboundRepository) {
				outboundRepoMock = mocksRepo.NewOutboundRepository(t)
				outboundRepoMock.On("Save", ctx, clientData, mock.Anything).Return(errors.New("Cannot Store to OutboundMessage"))
				return
			},
			loggerMock: func() (loggerMock *mocksProvider.ILogger) {
				loggerMock = mocksProvider.NewILogger(t)
				loggerMock.On("Infof", provider.AppLog, mock.Anything)
				loggerMock.On("Errorf", provider.AppLog, mock.Anything)
				return
			},
			want: want{
				err: errors.New("Cannot Store to OutboundMessage"),
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {

			outboundRepoMock := test.outboundRepoMock()

			drRepoMock := test.drRepoMock()

			metaPlatformMock := test.metaPlatformMock()
			loggerMock := test.loggerMock()

			service := NewQueueService(outboundRepoMock, drRepoMock, metaPlatformMock, loggerMock)

			msgBytes := []byte(test.msgString)

			err := service.SendMessage(ctx, msgBytes, test.clientData)

			assert.Equal(t, test.want.err, err)
			outboundRepoMock.AssertExpectations(t)
			drRepoMock.AssertExpectations(t)
			metaPlatformMock.AssertExpectations(t)
			loggerMock.AssertExpectations(t)
		})
	}
}
