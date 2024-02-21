package consumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"transmitter-artemis/entity"
	"transmitter-artemis/provider"
	"transmitter-artemis/service"

	"github.com/go-stomp/stomp/v3"
)

type QueueListener struct {
	conn         *stomp.Conn
	sub          *stomp.Subscription
	queueName    string
	clientData   entity.ClientData
	queueService service.QueueService
	logger       provider.ILogger
}

func NewQueueListener(conn *stomp.Conn, queueService service.QueueService, clientData entity.ClientData, logger provider.ILogger) *QueueListener {

	queueName := fmt.Sprintf("%s-msg-queue", clientData.ClientName)

	sub, err := conn.Subscribe(queueName, stomp.AckClientIndividual)
	if err != nil {
		logger.Errorf(provider.AppLog, "Failed to subscribe to queue:%s,error: %v", queueName, err)
	}
	logger.Infof(provider.AppLog, "Success to subscribe to queue:%s", queueName)
	return &QueueListener{
		conn:         conn,
		sub:          sub,
		queueName:    queueName,
		clientData:   clientData,
		queueService: queueService,
		logger:       logger,
	}
}

func (ql *QueueListener) Start(ctx context.Context) {
	// Membuat channel untuk menerima sinyal SIGINT atau SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg := <-ql.sub.C:

			ql.logger.Infof(provider.AmqLog, "From %s, Received message: %s", ql.queueName, string(msg.Body))

			ql.queueService.SendMessage(ctx, msg.Body, ql.clientData)

			err := ql.conn.Ack(msg)
			if err != nil {
				ql.logger.Errorf(provider.AmqLog, "Failed to acknowledge message: %v", err)
			}
		case <-stop:
			log.Println("Received signal to stop. Exiting...")
			ql.sub.Unsubscribe()
			return
		case <-ctx.Done():
			log.Println("Context is canceled. Exiting...")
			ql.sub.Unsubscribe()
			return
		}
	}
}
