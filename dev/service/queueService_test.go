package service

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"testing"
// 	"time"
// 	"transmitter-artemis/dto"
// 	"transmitter-artemis/entity"
// 	platform "transmitter-artemis/mocks/platform"
// 	util "transmitter-artemis/mocks/provider"
// 	mocks "transmitter-artemis/mocks/repository"
// 	"transmitter-artemis/provider"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestSendMessage(t *testing.T) {

// 	t.Run("Test Success SendMessage", func(t *testing.T) {
// 		mockLogger := util.NewILogger(t)
// 		mockDRRepo := mocks.NewDRRepository(t)
// 		mockOutboundRepo := mocks.NewOutboundRepository(t)
// 		mockMetaPlatform := platform.NewMetaClient(t)

// 		app := NewQueueService(mockOutboundRepo, mockDRRepo, mockMetaPlatform, mockLogger)
// 		ctx := context.Background()

// 		msgString := "message_id=1&to=6283872750005&type=text&text[preview_url]=false&text[body]=Hello, this is a test message"
// 		msgBytes := []byte(msgString)

// 		payload := dto.RequestToMeta{

// 			MessagingProduct: "whatsapp",
// 			RecipientType:    "individual",
// 			To:               "6283872750005",
// 			Type:             "text",
// 			Text: struct {
// 				PreviewURL bool   `json:"preview_url" bson:"preview_url"`
// 				Body       string `json:"body" bson:"body"`
// 			}{
// 				PreviewURL: false,
// 				Body:       "Hello, this is a test message",
// 			},
// 		}

// 		response := dto.ResponseFromMeta{
// 			MessagingProduct: "whatsapp",
// 			Contacts: []struct {
// 				Input string `json:"input" bson:"input"`
// 				WAID  string `json:"wa_id,omitempty" bson:"wa_id,omitempty"`
// 			}{
// 				{Input: "6283872750005", WAID: "6283872750005"},
// 			},
// 			Messages: []struct {
// 				ID string `json:"id"`
// 			}{
// 				{ID: "wamid.abc123"},
// 			},
// 		}

// 		clientData := entity.ClientData{
// 			ClientName:    "lino",
// 			Token:         "abc",
// 			PhoneNumberID: "123",
// 			WAHost:        "https://graph.facebook.com",
// 		}

// 		URL := fmt.Sprintf("%s/%s/messages", clientData.WAHost, clientData.PhoneNumberID)

// 		var statusCode int

// 		rightNow := time.Now()
// 		timeString := fmt.Sprintf("%d", rightNow.Unix())
// 		drMsg := fmt.Sprintf("message_id=1&wa_id=wamid.abc123&deliverystatus=sent&time=%v", timeString)

// 		mockMetaPlatform.On("SendRequestToMeta", ctx, URL, clientData.Token, payload).Return(response, statusCode, nil)

// 		mockDRRepo.On("Produce", ctx, clientData, drMsg).Return(nil)

// 		mockLogger.On("Infof", provider.AppLog, "Success Store to DR-MSG")

// 		mockOutboundRepo.On("Save", ctx, clientData, mock.Anything).Return(nil)

// 		mockLogger.On("Infof", provider.AppLog, "Success Store Data to OutboundMessage")
// 		err := app.SendMessage(ctx, msgBytes, clientData)

// 		assert.NoError(t, err)
// 		mockDRRepo.AssertExpectations(t)
// 		mockOutboundRepo.AssertExpectations(t)
// 		mockLogger.AssertExpectations(t)
// 		mockMetaPlatform.AssertExpectations(t)
// 	})

// 	t.Run("Test Failed SendMessage", func(t *testing.T) {
// 		mockLogger := util.NewILogger(t)
// 		mockDRRepo := mocks.NewDRRepository(t)
// 		mockOutboundRepo := mocks.NewOutboundRepository(t)
// 		mockMetaPlatform := platform.NewMetaClient(t)

// 		app := NewQueueService(mockOutboundRepo, mockDRRepo, mockMetaPlatform, mockLogger)
// 		ctx := context.Background()

// 		msgString := "message_id=1&to=6283872750005&type=text&text[preview_url]=false&text[body]=Hello, this is a test message"
// 		msgBytes := []byte(msgString)

// 		payload := dto.RequestToMeta{

// 			MessagingProduct: "whatsapp",
// 			RecipientType:    "individual",
// 			To:               "6283872750005",
// 			Type:             "text",
// 			Text: struct {
// 				PreviewURL bool   `json:"preview_url" bson:"preview_url"`
// 				Body       string `json:"body" bson:"body"`
// 			}{
// 				PreviewURL: false,
// 				Body:       "Hello, this is a test message",
// 			},
// 		}

// 		errorRes := dto.ErrorRes{
// 			Message: "Error message",
// 		}

// 		response := dto.ResponseFromMeta{
// 			Error: &errorRes,
// 		}

// 		clientData := entity.ClientData{
// 			ClientName:    "lino",
// 			Token:         "abc",
// 			PhoneNumberID: "123",
// 			WAHost:        "https://graph.facebook.com",
// 		}

// 		URL := fmt.Sprintf("%s/%s/messages", clientData.WAHost, clientData.PhoneNumberID)

// 		var statusCode int

// 		mockMetaPlatform.On("SendRequestToMeta", ctx, URL, clientData.Token, payload).Return(response, statusCode, nil)

// 		mockDRRepo.On("Produce", ctx, clientData, mock.Anything).Return(nil)

// 		mockLogger.On("Infof", provider.AppLog, "Success Store to DR-MSG")

// 		mockOutboundRepo.On("Save", ctx, clientData, mock.Anything).Return(nil)

// 		mockLogger.On("Infof", provider.AppLog, mock.Anything)
// 		err := app.SendMessage(ctx, msgBytes, clientData)

// 		assert.NoError(t, err)
// 		mockDRRepo.AssertExpectations(t)
// 		mockOutboundRepo.AssertExpectations(t)
// 		mockLogger.AssertExpectations(t)
// 		mockMetaPlatform.AssertExpectations(t)
// 	})

// 	t.Run("Test Cannot Parse queue to request body", func(t *testing.T) {
// 		mockLogger := util.NewILogger(t)
// 		mockDRRepo := mocks.NewDRRepository(t)
// 		mockOutboundRepo := mocks.NewOutboundRepository(t)
// 		mockMetaPlatform := platform.NewMetaClient(t)

// 		app := NewQueueService(mockOutboundRepo, mockDRRepo, mockMetaPlatform, mockLogger)
// 		clientData := entity.ClientData{
// 			ClientName:    "lino",
// 			Token:         "abc",
// 			PhoneNumberID: "123",
// 			WAHost:        "https://graph.facebook.com",
// 		}
// 		msgString := "to=6283872750005&type;=text" // <-- invalid input
// 		msgBytes := []byte(msgString)
// 		ctx := context.Background()
// 		mockLogger.On("Errorf", provider.AppLog, mock.Anything)
// 		err := app.SendMessage(ctx, msgBytes, clientData)
// 		assert.Error(t, err)
// 		mockLogger.AssertExpectations(t)

// 	})

// 	t.Run("Test Cannot Save DR queue to Artemis", func(t *testing.T) {
// 		mockLogger := util.NewILogger(t)
// 		mockDRRepo := mocks.NewDRRepository(t)
// 		mockOutboundRepo := mocks.NewOutboundRepository(t)
// 		mockMetaPlatform := platform.NewMetaClient(t)

// 		app := NewQueueService(mockOutboundRepo, mockDRRepo, mockMetaPlatform, mockLogger)
// 		clientData := entity.ClientData{
// 			ClientName:    "lino",
// 			Token:         "abc",
// 			PhoneNumberID: "123",
// 			WAHost:        "https://graph.facebook.com",
// 		}
// 		ctx := context.Background()

// 		msgString := "invalid_query_string"
// 		msgBytes := []byte(msgString)
// 		response := dto.ResponseFromMeta{
// 			MessagingProduct: "whatsapp",
// 			Contacts: []struct {
// 				Input string `json:"input" bson:"input"`
// 				WAID  string `json:"wa_id,omitempty" bson:"wa_id,omitempty"`
// 			}{
// 				{Input: "6283872750005", WAID: "6283872750005"},
// 			},
// 			Messages: []struct {
// 				ID string `json:"id"`
// 			}{
// 				{ID: "wamid.abc123"},
// 			},
// 		}

// 		err := errors.New("Cannot Save to DR-queue Artemis")

// 		mockMetaPlatform.On("SendRequestToMeta", ctx, mock.Anything, mock.Anything, mock.Anything).Return(response, 200, nil)

// 		mockDRRepo.On("Produce", ctx, clientData, mock.Anything).Return(err)
// 		mockLogger.On("Errorf", provider.AppLog, mock.Anything)
// 		err = app.SendMessage(ctx, msgBytes, clientData)
// 		t.Log(err)
// 		assert.Error(t, err)
// 		mockDRRepo.AssertExpectations(t)
// 		mockOutboundRepo.AssertExpectations(t)
// 		mockLogger.AssertExpectations(t)
// 		mockMetaPlatform.AssertExpectations(t)
// 	})

// 	t.Run("Test Cant Send request to meta", func(t *testing.T) {
// 		mockLogger := util.NewILogger(t)
// 		mockDRRepo := mocks.NewDRRepository(t)
// 		mockOutboundRepo := mocks.NewOutboundRepository(t)
// 		mockMetaPlatform := platform.NewMetaClient(t)

// 		app := NewQueueService(mockOutboundRepo, mockDRRepo, mockMetaPlatform, mockLogger)
// 		clientData := entity.ClientData{
// 			ClientName:    "lino",
// 			Token:         "abc",
// 			PhoneNumberID: "123",
// 			WAHost:        "https://graph.facebook.com",
// 		}
// 		ctx := context.Background()

// 		msgString := "invalid_query_string"
// 		msgBytes := []byte(msgString)

// 		mockMetaPlatform.On("SendRequestToMeta", ctx, mock.Anything, mock.Anything, mock.Anything).Return(dto.ResponseFromMeta{}, 500, errors.New("Cannot Send Request"))

// 		mockLogger.On("Errorf", provider.AppLog, mock.Anything)

// 		err := app.SendMessage(ctx, msgBytes, clientData)
// 		t.Log(err)
// 		assert.Error(t, err)
// 		mockDRRepo.AssertExpectations(t)
// 		mockOutboundRepo.AssertExpectations(t)
// 		mockLogger.AssertExpectations(t)
// 		mockMetaPlatform.AssertExpectations(t)
// 	})

// 	t.Run("Test Cannot Store msg to Mongo", func(t *testing.T) {
// 		mockLogger := util.NewILogger(t)
// 		mockDRRepo := mocks.NewDRRepository(t)
// 		mockOutboundRepo := mocks.NewOutboundRepository(t)
// 		mockMetaPlatform := platform.NewMetaClient(t)

// 		app := NewQueueService(mockOutboundRepo, mockDRRepo, mockMetaPlatform, mockLogger)
// 		clientData := entity.ClientData{
// 			ClientName:    "lino",
// 			Token:         "abc",
// 			PhoneNumberID: "123",
// 			WAHost:        "https://graph.facebook.com",
// 		}
// 		ctx := context.Background()

// 		msgString := "invalid_query_string"
// 		msgBytes := []byte(msgString)
// 		response := dto.ResponseFromMeta{
// 			MessagingProduct: "whatsapp",
// 			Contacts: []struct {
// 				Input string `json:"input" bson:"input"`
// 				WAID  string `json:"wa_id,omitempty" bson:"wa_id,omitempty"`
// 			}{
// 				{Input: "6283872750005", WAID: "6283872750005"},
// 			},
// 			Messages: []struct {
// 				ID string `json:"id"`
// 			}{
// 				{ID: "wamid.abc123"},
// 			},
// 		}

// 		mockMetaPlatform.On("SendRequestToMeta", ctx, mock.Anything, mock.Anything, mock.Anything).Return(response, 200, nil)

// 		mockDRRepo.On("Produce", ctx, clientData, mock.Anything).Return(nil)
// 		mockLogger.On("Infof", provider.AppLog, mock.Anything)
// 		mockOutboundRepo.On("Save", ctx, clientData, mock.Anything).Return(errors.New("Cannot Store to OutboundMessage"))
// 		mockLogger.On("Errorf", provider.AppLog, mock.Anything)

// 		err := app.SendMessage(ctx, msgBytes, clientData)

// 		assert.Error(t, err)
// 		mockDRRepo.AssertExpectations(t)
// 		mockOutboundRepo.AssertExpectations(t)
// 		mockLogger.AssertExpectations(t)
// 		mockMetaPlatform.AssertExpectations(t)
// 	})

// }
