package testing

import (
	"fmt"
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

func (suite *testSvc) TestIntegrationTransmitter() {

	loc, _ := time.LoadLocation("Asia/Jakarta")
	timeNow := time.Date(2024, 02, 23, 18, 0, 0, 0, loc)

	t := suite.T()
	type want struct {
		outbound    entity.OutboundMessage
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
	}{{
		name:      "Test Success Message",
		queueBody: "message_id=1&to=valid_phone_number_id&type=text&text[preview_url]=false&text[body]=contoh Pesan",
		clientData: entity.ClientData{
			ClientName:    "test-client-a",
			Token:         "valid_token",
			PhoneNumberID: "valid_phone_number_id",
			WAHost:        suite.metaServer.URL,
		},
		msgID: "1",
		want: want{
			outbound: entity.OutboundMessage{
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
					"message_id=1&wa_id=wamid.123&deliverystatus=sent&time=112233",
				},
				CreatedAt: timeNow,
				UpdatedAt: timeNow,
			},
			dr:          "message_id=1&wa_id=wamid.123&deliverystatus=sent&time=112233",
			checkDB:     true,
			isProduceDR: true,
		},
	}}

	// var drData []string
	_ = suite.validateClientData(t)
	for no, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("Test Case:", no+1)

			queueName := fmt.Sprintf("%s-msg-queue", tt.clientData.ClientName)
			time.Sleep(2 * time.Second)
			suite.insertQueue(t, tt.queueBody, queueName)
			time.Sleep(2 * time.Second)
			// suite.validateQueueMessage(t, tt.queueBody, queueName)

			// if tt.want.isProduceDR {
			// 	drData = append(drData, tt.want.dr)
			// }

			// if tt.want.checkDB {
			// 	time.Sleep(2 * time.Second)
			// 	actualOutbound := suite.validateDB(t, tt.msgID, tt.clientData, tt.want.outbound)
			// 	assert.Equal(t, tt.want.outbound, actualOutbound)
			// }

		})
	}
	// time.Sleep(1 * time.Second)
	// suite.validateDRQueue(t, drData)

}

func (s *testSvc) insertQueue(t *testing.T, queueData string, queueName string) {

	contentType := "text/plain"
	err := s.artemis.Send(queueName, contentType, []byte(queueData), stomp.SendOpt.Header("destination-type", "ANYCAST"), stomp.SendOpt.Header("persistent", "true"))
	// fmt.Println("SEND MSG TO QUEUE DONE")

	// fmt.Println(queueData)
	fmt.Println("INSERT To queue name : ", queueName)
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
	fmt.Println("expected queue length : ", len(expected))
	queueName := "test-client-a-dr-msg"

	subs, err := s.artemis.Subscribe(queueName, stomp.AckAuto)
	defer subs.Unsubscribe()
	assert.NoError(t, err)
	assert.NotNil(t, subs)

	// done := make(chan struct{})
	// defer close(done)

	i := 0

	// Menangani sinyal SIGINT (ctrl+c) untuk membatalkan konteks
	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for msgQueue := range subs.C {
		fmt.Println("Msg ke-", i+1)
		fmt.Println("Actual message: ", expected[i])
		actualMsg := string(msgQueue.Body)
		fmt.Println("Received message: ", actualMsg)

		assert.Equal(t, expected[i], actualMsg)

		i++

		if i == len(expected) {
			break
		}

		// Memeriksa apakah konteks telah dibatalkan
		select {
		case <-s.ctx.Done():
			fmt.Println("MEMBERHENTIKAN TESTING ...")
			return // Keluar dari loop jika konteks telah dibatalkan
		default:
		}
	}

	err = subs.Unsubscribe()
	assert.NoError(t, err)
}

func (s *testSvc) validateDB(t *testing.T, msgID string, clientData entity.ClientData, expectedOutbound entity.OutboundMessage) entity.OutboundMessage {
	outboundMessage := entity.OutboundMessage{}

	collName := fmt.Sprintf("%s-outbound-msg", clientData.ClientName)
	fmt.Println("Coll name =", collName)
	fmt.Println("Msg ID =", msgID)
	collection := s.mongoClient.Database(config.Configuration.MongoDB.Database).Collection(collName)
	filter := bson.D{
		{
			Key:   "message_id",
			Value: msgID,
		},
	}
	err := collection.FindOne(s.ctx, filter).Decode(&outboundMessage)
	fmt.Println("DATA YANG PADA MONGODB:", outboundMessage)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutbound, outboundMessage)
	return outboundMessage
}

func (s *testSvc) validateClientData(t *testing.T) (clientData entity.ClientData) {

	// collName := fmt.Sprintf("%s-outbound-msg", clientData.ClientName)
	collName := "client-info"
	collection := s.mongoClient.Database(config.Configuration.MongoDB.Database).Collection(collName)
	filter := bson.D{
		{
			Key:   "client_name",
			Value: "test-client-a",
		},
	}
	err := collection.FindOne(s.ctx, filter).Decode(&clientData)
	fmt.Println("DATA CLIENT YANG PADA MONGODB:", clientData)
	assert.NoError(t, err)

	return clientData
}
