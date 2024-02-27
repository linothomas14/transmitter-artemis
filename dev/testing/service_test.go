package testing

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"
	"transmitter-artemis/config"
	"transmitter-artemis/dto"
	"transmitter-artemis/entity"

	"github.com/go-stomp/stomp/v3"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

type drAndClientData struct {
	drMsg      []byte
	clientData entity.ClientData
}

type outboundMessageTestable struct {
	MessageID        string               `json:"message_id" bson:"message_id"`
	WAID             string               `json:"wa_id,omitempty" bson:"wa_id,omitempty"`
	To               string               `json:"to" bson:"to"`
	OriginalRequest  string               `json:"original_request" bson:"original_request"`
	Request          dto.RequestToMeta    `json:"request"`
	OriginalResponse dto.ResponseFromMeta `json:"original_response" bson:"original_response"`
	DeliveryReport   []string             `json:"delivery_report" bson:"delivery_report"`
}

func (suite *testSvc) TestIntegrationTransmitter() {

	t := suite.T()
	type want struct {
		outbound    outboundMessageTestable
		dr          string
		isProduceDR bool
		checkDB     bool
	}

	tests := []struct {
		name       string
		clientData entity.ClientData
		queueBody  string
		msgID      string
		want       want
	}{
		{
			msgID:     "1",
			name:      "Test Happy Success Message",
			queueBody: "message_id=1&to=valid_phone_number&type=text&text[preview_url]=false&text[body]=contoh Pesan",
			clientData: entity.ClientData{
				ClientName:    ConstClientName,
				Token:         "valid_token",
				PhoneNumberID: "valid_phone_number_id",
				WAHost:        suite.metaServer.URL,
			},
			want: want{
				outbound: outboundMessageTestable{
					MessageID:       "1",
					WAID:            "wamid.123",
					To:              "valid_phone_number",
					OriginalRequest: "message_id=1&to=valid_phone_number&type=text&text[preview_url]=false&text[body]=contoh Pesan",
					Request: dto.RequestToMeta{
						MessagingProduct: "whatsapp",
						RecipientType:    "individual",
						To:               "valid_phone_number",
						Type:             "text",
						Text: struct {
							PreviewURL bool   `json:"preview_url" bson:"preview_url"`
							Body       string `json:"body" bson:"body"`
						}{
							PreviewURL: false,
							Body:       "contoh Pesan",
						},
					},
					OriginalResponse: dto.ResponseFromMeta{
						MessagingProduct: "whatsapp",
						Contacts: []struct {
							Input string `json:"input" bson:"input"`
							WAID  string `json:"wa_id,omitempty" bson:"wa_id,omitempty"`
						}{
							{
								Input: "valid_phone_number",
								WAID:  "valid_phone_number",
							},
						},
						Messages: []struct {
							ID string `json:"id"`
						}{
							{
								ID: "wamid.123",
							},
						},
					},
					DeliveryReport: []string{
						"message_id=1&wa_id=wamid.123&deliverystatus=sent",
					},
				},
				dr:          "message_id=1&wa_id=wamid.123&deliverystatus=sent",
				checkDB:     true,
				isProduceDR: true,
			},
		},
		{
			msgID:     "2",
			name:      "Test Phone Number not registered yet",
			queueBody: "message_id=2&to=invalid_phone_number&type=text&text[preview_url]=false&text[body]=contoh Pesan",
			clientData: entity.ClientData{
				ClientName:    ConstClientName,
				Token:         "valid_token",
				PhoneNumberID: "valid_phone_number_id",
				WAHost:        suite.metaServer.URL,
			},
			want: want{
				outbound: outboundMessageTestable{
					MessageID:       "2",
					To:              "invalid_phone_number",
					OriginalRequest: "message_id=2&to=invalid_phone_number&type=text&text[preview_url]=false&text[body]=contoh Pesan",
					Request: dto.RequestToMeta{
						MessagingProduct: "whatsapp",
						RecipientType:    "individual",
						To:               "invalid_phone_number",
						Type:             "text",
						Text: struct {
							PreviewURL bool   `json:"preview_url" bson:"preview_url"`
							Body       string `json:"body" bson:"body"`
						}{
							PreviewURL: false,
							Body:       "contoh Pesan",
						},
					},
					OriginalResponse: dto.ResponseFromMeta{

						Error: &dto.ErrorRes{
							Message: "(#131030) Recipient phone number not in allowed list",
							Type:    "OAuthException",
							Code:    131030,
							ErrorData: map[string]interface{}{
								"messaging_product": "whatsapp",
								"details":           "Recipient phone number not in allowed list: Add recipient phone number…",
							},
							FbTraceID: "INVALID PHONE NUMBER",
						},
					},
					DeliveryReport: []string{
						"message_id=2&deliverystatus=failed&error[code]=131030&error[detail]=(#131030) Recipient phone number not in allowed list",
					},
				},
				dr:          "message_id=2&deliverystatus=failed&error[code]=131030&error[detail]=(#131030) Recipient phone number not in allowed list",
				checkDB:     true,
				isProduceDR: true,
			},
		},
		{
			msgID:     "3",
			name:      "Test text body was empty",
			queueBody: "message_id=3&to=valid_phone_number&type=text&text[preview_url]=false&text[body]=",
			clientData: entity.ClientData{
				ClientName:    ConstClientName,
				Token:         "valid_token",
				PhoneNumberID: "valid_phone_number_id",
				WAHost:        suite.metaServer.URL,
			},
			want: want{
				outbound: outboundMessageTestable{
					MessageID:       "3",
					To:              "valid_phone_number",
					OriginalRequest: "message_id=3&to=valid_phone_number&type=text&text[preview_url]=false&text[body]=",
					Request: dto.RequestToMeta{
						MessagingProduct: "whatsapp",
						RecipientType:    "individual",
						To:               "valid_phone_number",
						Type:             "text",
						Text: struct {
							PreviewURL bool   `json:"preview_url" bson:"preview_url"`
							Body       string `json:"body" bson:"body"`
						}{
							PreviewURL: false,
							Body:       "",
						},
					},
					OriginalResponse: dto.ResponseFromMeta{

						Error: &dto.ErrorRes{
							Message: "(#100) The parameter text['body'] is required.",
							Type:    "OAuthException",
							Code:    100,

							FbTraceID: "BODY MSG WAS EMPTY",
						},
					},
					DeliveryReport: []string{
						"message_id=3&deliverystatus=failed&error[code]=100&error[detail]=(#100) The parameter text['body'] is required.",
					},
				},
				dr:          "message_id=3&deliverystatus=failed&error[code]=100&error[detail]=(#100) The parameter text['body'] is required.",
				checkDB:     true,
				isProduceDR: true,
			},
		},
		{
			msgID:     "4",
			name:      "Test Type not text",
			queueBody: "message_id=4&to=valid_phone_number&type=not-text&text[preview_url]=false&text[body]=contoh Pesan",
			clientData: entity.ClientData{
				ClientName:    ConstClientName,
				Token:         "valid_token",
				PhoneNumberID: "valid_phone_number_id",
				WAHost:        suite.metaServer.URL,
			},
			want: want{
				outbound: outboundMessageTestable{
					MessageID:       "4",
					To:              "valid_phone_number",
					OriginalRequest: "message_id=4&to=valid_phone_number&type=not-text&text[preview_url]=false&text[body]=contoh Pesan",
					Request: dto.RequestToMeta{
						MessagingProduct: "whatsapp",
						RecipientType:    "individual",
						To:               "valid_phone_number",
						Type:             "not-text",
						Text: struct {
							PreviewURL bool   `json:"preview_url" bson:"preview_url"`
							Body       string `json:"body" bson:"body"`
						}{
							PreviewURL: false,
							Body:       "contoh Pesan",
						},
					},
					OriginalResponse: dto.ResponseFromMeta{
						Error: &dto.ErrorRes{
							Message:   "(#100) Param type must be one of {AUDIO, CONTACTS, DOCUMENT, IMAGE, INTERACTIVE, LINK_PREVIEW, LOCATION, REACTION, STICKER, TEMPLATE, TEXT, VIDEO} - got \"not-text\".",
							Type:      "OAuthException",
							Code:      100,
							FbTraceID: "Type must be Text",
						},
					},
					DeliveryReport: []string{
						"message_id=4&deliverystatus=failed&error[code]=100&error[detail]=(#100) Param type must be one of {AUDIO, CONTACTS, DOCUMENT, IMAGE, INTERACTIVE, LINK_PREVIEW, LOCATION, REACTION, STICKER, TEMPLATE, TEXT, VIDEO} - got \"not-text\".",
					},
				},
				dr:          "message_id=4&deliverystatus=failed&error[code]=100&error[detail]=(#100) Param type must be one of {AUDIO, CONTACTS, DOCUMENT, IMAGE, INTERACTIVE, LINK_PREVIEW, LOCATION, REACTION, STICKER, TEMPLATE, TEXT, VIDEO} - got \"not-text\".",
				checkDB:     true,
				isProduceDR: true,
			},
		},
		{
			msgID:     "5",
			name:      "Test no \"To\" Field in queueBody",
			queueBody: "message_id=5&type=text&text[preview_url]=false&text[body]=contoh Pesan",
			clientData: entity.ClientData{
				ClientName:    ConstClientName,
				Token:         "valid_token",
				PhoneNumberID: "valid_phone_number_id",
				WAHost:        suite.metaServer.URL,
			},
			want: want{
				outbound: outboundMessageTestable{
					MessageID:       "5",
					To:              "",
					OriginalRequest: "message_id=5&type=text&text[preview_url]=false&text[body]=contoh Pesan",
					Request: dto.RequestToMeta{
						MessagingProduct: "whatsapp",
						RecipientType:    "individual",
						To:               "",
						Type:             "text",
						Text: struct {
							PreviewURL bool   `json:"preview_url" bson:"preview_url"`
							Body       string `json:"body" bson:"body"`
						}{
							PreviewURL: false,
							Body:       "contoh Pesan",
						},
					},
					OriginalResponse: dto.ResponseFromMeta{
						Error: &dto.ErrorRes{
							Message:   "The parameter to is required.",
							Type:      "OAuthException",
							Code:      100,
							FbTraceID: "NO \"TO\" FIELD",
						},
					},
					DeliveryReport: []string{
						"message_id=5&deliverystatus=failed&error[code]=100&error[detail]=The parameter to is required.",
					},
				},
				dr:          "message_id=5&deliverystatus=failed&error[code]=100&error[detail]=The parameter to is required.",
				checkDB:     true,
				isProduceDR: true,
			},
		},
		{
			msgID:     "6",
			name:      "Test Invalid Token",
			queueBody: "message_id=6&to=valid_phone_number&type=text&text[preview_url]=false&text[body]=contoh Pesan",
			clientData: entity.ClientData{
				ClientName:    "client-2",
				Token:         "invalid_token",
				PhoneNumberID: "valid_phone_number_id",
				WAHost:        suite.metaServer.URL,
			},
			want: want{
				outbound: outboundMessageTestable{
					MessageID:       "6",
					To:              "valid_phone_number",
					OriginalRequest: "message_id=6&to=valid_phone_number&type=text&text[preview_url]=false&text[body]=contoh Pesan",
					Request: dto.RequestToMeta{
						MessagingProduct: "whatsapp",
						RecipientType:    "individual",
						To:               "valid_phone_number",
						Type:             "text",
						Text: struct {
							PreviewURL bool   `json:"preview_url" bson:"preview_url"`
							Body       string `json:"body" bson:"body"`
						}{
							PreviewURL: false,
							Body:       "contoh Pesan",
						},
					},
					OriginalResponse: dto.ResponseFromMeta{

						Error: &dto.ErrorRes{
							Message: "Invalid OAuth access token - Cannot parse access token",
							Type:    "OAuthException",
							Code:    190,
							ErrorData: map[string]interface{}{
								"messaging_product": "whatsapp",
								"details":           "Invalid OAuth access token - Cannot parse access token",
							},
							FbTraceID: "INVALID TOKEN",
						},
					},
					DeliveryReport: []string{
						"message_id=6&deliverystatus=failed&error[code]=190&error[detail]=Invalid OAuth access token - Cannot parse access token",
					},
				},
				dr:          "message_id=6&deliverystatus=failed&error[code]=190&error[detail]=Invalid OAuth access token - Cannot parse access token",
				checkDB:     true,
				isProduceDR: true,
			},
		},
	}

	var drData []string
	suite.validateClientData(t)
	for no, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("Test Case:", no+1)

			queueName := fmt.Sprintf("%s-msg-queue", tt.clientData.ClientName)
			suite.insertQueue(t, tt.queueBody, queueName)

			fmt.Println("Inserting to queue : ", queueName)
			fmt.Println("The data : ", tt.queueBody)

			if tt.want.isProduceDR {
				drData = append(drData, tt.want.dr)
			}

			if tt.want.checkDB {
				time.Sleep(1 * time.Second)
				suite.validateDB(t, tt.msgID, tt.clientData, tt.want.outbound)
			}
		})
	}
	suite.validateDRQueue(t, drData)

}

func (s *testSvc) insertQueue(t *testing.T, queueData string, queueName string) {

	contentType := "text/plain"
	err := s.artemis.Send(queueName, contentType, []byte(queueData), stomp.SendOpt.Header("destination-type", "ANYCAST"), stomp.SendOpt.Header("persistent", "true"))
	fmt.Println("SEND MSG TO", queueName, "SUCCESS")
	assert.NoError(t, err)
}
func (s *testSvc) validateQueueMessage(t *testing.T, queueData string, queueName string) {

	subs, err := s.artemis.Subscribe(queueName, stomp.AckClientIndividual)
	defer subs.Unsubscribe()
	fmt.Println("Subribing to queue name : ", queueName)
	fmt.Println("LISTENING...")
	assert.NoError(t, err)
	assert.NotNil(t, subs)

	for msgQueue := range subs.C {

		actualMsg := string(msgQueue.Body)
		fmt.Println("Received message: ", actualMsg)
		break
	}

	assert.NoError(t, err)

}
func (s *testSvc) validateDRQueue(t *testing.T, expected []string) {
	fmt.Println("Expected queue length : ", len(expected))

	testForInvalidClient := 1
	testForValidClient := len(expected) - testForInvalidClient

	// FOR VALID CLIENT TOKEN
	queueNameClient1 := fmt.Sprintf("%s-dr-msg", ConstClientName)
	subs1, err := s.artemis.Subscribe(queueNameClient1, stomp.AckAuto)
	defer subs1.Unsubscribe()

	// FOR INVALID CLIENT TOKEN
	queueNameClient2 := fmt.Sprintf("%s-dr-msg", "client-2")
	subs2, err := s.artemis.Subscribe(queueNameClient2, stomp.AckAuto)
	defer subs2.Unsubscribe()
	assert.NoError(t, err)
	assert.NotNil(t, subs1)
	assert.NoError(t, err)
	assert.NotNil(t, subs2)

	i := 0

	for msgQueue := range subs1.C {
		fmt.Println("Msg ke-", i+1)
		fmt.Println("Expected message on Artemis: ", expected[i])
		actualMsg := string(msgQueue.Body)
		fmt.Println("Actual message on Artemis: ", actualMsg)

		msgWithoutTime := strings.Split(actualMsg, "&time=")
		fmt.Println("Converted message to : ", msgWithoutTime[0])
		assert.Equal(t, expected[i], msgWithoutTime[0])

		i++

		// TOTAL TESTCASE WITH VALID CLIENT
		if i == testForValidClient {
			break
		}
	}

	for msgQueue := range subs2.C {
		fmt.Println("Msg ke-", i+1)
		fmt.Println("Expected message on Artemis: ", expected[i])
		actualMsg := string(msgQueue.Body)
		fmt.Println("Actual message on Artemis: ", actualMsg)

		msgWithoutTime := strings.Split(actualMsg, "&time=")
		fmt.Println("Converted message to : ", msgWithoutTime[0])
		assert.Equal(t, expected[i], msgWithoutTime[0])

		i++

		// TOTAL TESTCASE WITH VALID CLIENT
		if i == len(expected) {
			break
		}

	}

	err = subs1.Unsubscribe()
	assert.NoError(t, err)
	err = subs2.Unsubscribe()
	assert.NoError(t, err)
}

func (s *testSvc) validateDB(t *testing.T, msgID string, clientData entity.ClientData, expectedOutbound outboundMessageTestable) {
	actualOutbound := entity.OutboundMessage{}

	collName := fmt.Sprintf("%s-outbound-msg", clientData.ClientName)

	collection := s.mongoClient.Database(config.Configuration.MongoDB.Database).Collection(collName)
	filter := bson.D{
		{
			Key:   "message_id",
			Value: msgID,
		},
	}
	err := collection.FindOne(s.ctx, filter).Decode(&actualOutbound)
	fmt.Println("OutboundMessage on MongoDB : ", actualOutbound)
	assert.NoError(t, err)
	outboundMessage := convertToTestableFormat(actualOutbound)
	assert.Equal(t, expectedOutbound, outboundMessage)
}

func (s *testSvc) validateClientData(t *testing.T) {

	var clients []entity.ClientData

	collName := "client-info"

	collection := s.mongoClient.Database(config.Configuration.MongoDB.Database).Collection(collName)

	cursor, err := collection.Find(context.Background(), bson.M{})

	assert.NoError(t, err)

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var result entity.ClientData
		if err := cursor.Decode(&result); err != nil {
			log.Println(err)
			continue
		}
		clients = append(clients, result)
	}

	assert.NoError(t, err)

	fmt.Println("DATA CLIENT YANG PADA MONGODB:", clients)
	assert.NoError(t, err)

}

func convertToTestableFormat(outboundMessage entity.OutboundMessage) outboundMessageTestable {

	var deliveryReport []string

	drWithOutTime := strings.Split(outboundMessage.DeliveryReport[0], "&time=")

	deliveryReport = append(deliveryReport, drWithOutTime[0])

	return outboundMessageTestable{
		MessageID:        outboundMessage.MessageID,
		WAID:             outboundMessage.WAID,
		To:               outboundMessage.To,
		OriginalRequest:  outboundMessage.OriginalRequest,
		Request:          outboundMessage.Request,
		OriginalResponse: outboundMessage.OriginalResponse,
		DeliveryReport:   deliveryReport,
	}
}
