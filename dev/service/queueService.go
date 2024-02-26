package service

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"
	"transmitter-artemis/dto"
	"transmitter-artemis/entity"
	"transmitter-artemis/platform"
	"transmitter-artemis/provider"
	"transmitter-artemis/repository"
)

type QueueService interface {
	SendMessage(ctx context.Context, msg []byte, clientData entity.ClientData) error
}

type queueService struct {
	outboundRepo repository.OutboundRepository
	drRepo       repository.DRRepository
	metaPlatform platform.MetaClient
	logger       provider.ILogger
}

func NewQueueService(outboundRepo repository.OutboundRepository, drRepo repository.DRRepository, metaPlatform platform.MetaClient, logger provider.ILogger) *queueService {
	return &queueService{
		outboundRepo: outboundRepo,
		drRepo:       drRepo,
		metaPlatform: metaPlatform,
		logger:       logger,
	}
}

func (qs *queueService) SendMessage(ctx context.Context, queueData []byte, clientData entity.ClientData) (err error) {

	var responseFromMeta dto.ResponseFromMeta
	// from queue to json
	msgRequest, message_id, err := TransformToRequestBody(queueData)

	if err != nil {
		qs.logger.Errorf(provider.AppLog, fmt.Sprintf("Cant Parse Queue to Req Body, %v", err.Error()))
		return
	}

	URL := fmt.Sprintf("%s/%s/messages", clientData.WAHost, clientData.PhoneNumberID)

	responseFromMeta, _, err = qs.metaPlatform.SendRequestToMeta(ctx, URL, clientData.Token, msgRequest)
	if err != nil {
		qs.logger.Errorf(provider.AppLog, "Cant Send Request from Meta")
		return
	}

	drMsg := FormatResponseToQueue(responseFromMeta, message_id)

	// Store to DR xx-dr-msg (Artemis)
	err = qs.drRepo.Produce(ctx, clientData, drMsg)
	if err != nil {
		qs.logger.Errorf(provider.AppLog, "Cannot Save to DR-queue Artemis")
		return
	}
	qs.logger.Infof(provider.AppLog, "Success Store to DR-MSG")

	outboundMessage := FormatDataToOutboundMessage(queueData, msgRequest, responseFromMeta, drMsg)

	// Store to collection xx-outbound-msg (MongoDB)
	err = qs.outboundRepo.Save(ctx, clientData, outboundMessage)
	if err != nil {
		qs.logger.Errorf(provider.AppLog, "Cannot Store to OutboundMessage")
		return
	}
	qs.logger.Infof(provider.AppLog, "Success Store Data to OutboundMessage")
	return nil
}

func TransformToRequestBody(msg []byte) (dto.RequestToMeta, string, error) {
	var queue dto.RequestToMeta
	var message_id string

	// Parse query string
	msgString := string(msg)
	values, err := url.ParseQuery(msgString)
	if err != nil {
		return dto.RequestToMeta{}, "", err
	}
	// Construct Queue from query parameters
	for key, val := range values {
		switch key {
		case "message_id":
			message_id = val[0]
		case "to":
			queue.To = val[0]
		case "type":
			queue.Type = val[0]
		case "text[preview_url]":
			queue.Text.PreviewURL, _ = strconv.ParseBool(val[0])
		case "text[body]":
			queue.Text.Body, _ = url.QueryUnescape(val[0])
		}
	}

	queue.MessagingProduct = "whatsapp"
	queue.RecipientType = "individual"

	return queue, message_id, nil
}

func FormatResponseToQueue(data dto.ResponseFromMeta, message_id string) string {
	timeNow := fmt.Sprintf("%d", time.Now().Unix())

	if len(data.Messages) != 0 {
		drStatus := "sent"
		queueString := fmt.Sprintf("message_id=%s&wa_id=%s&deliverystatus=%s&time=%v", message_id, data.Messages[0].ID, drStatus, timeNow)
		return queueString
	} else {
		drStatus := "failed"
		queueString := fmt.Sprintf("message_id=%s&deliverystatus=%s&error[code]=%d&error[detail]=%v&time=%v", message_id, drStatus, data.Error.Code, data.Error.Message, timeNow)
		return queueString
	}
}

func FormatDataToOutboundMessage(queueData []byte, request dto.RequestToMeta, response dto.ResponseFromMeta, dr string) entity.OutboundMessage {

	var outboundMessage entity.OutboundMessage

	valuesDR, _ := url.ParseQuery(dr)

	wa_id := valuesDR.Get("wa_id")

	if wa_id != "" {
		outboundMessage.WAID = wa_id
	}

	queueData_string := string(queueData)
	valuesQueueData, _ := url.ParseQuery(queueData_string)

	to := valuesQueueData.Get("to")
	msg_id := valuesQueueData.Get("message_id")

	outboundMessage.To = to
	outboundMessage.MessageID = msg_id
	outboundMessage.OriginalRequest = queueData_string
	outboundMessage.Request = request
	outboundMessage.OriginalResponse = response
	outboundMessage.DeliveryReport = append(outboundMessage.DeliveryReport, dr)
	outboundMessage.CreatedAt = time.Now()
	outboundMessage.UpdatedAt = time.Now()

	return outboundMessage
}
