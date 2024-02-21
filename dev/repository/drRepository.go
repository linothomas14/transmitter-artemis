package repository

import (
	"context"
	"fmt"
	"transmitter-artemis/entity"

	"github.com/go-stomp/stomp/v3"
)

type DRRepository interface {
	Produce(ctx context.Context, clientData entity.ClientData, drMsg string) error
}

type drRepository struct {
	conn *stomp.Conn
}

func NewDRRepository(conn *stomp.Conn) *drRepository {
	return &drRepository{
		conn: conn,
	}
}

func (drRep *drRepository) Produce(ctx context.Context, clientData entity.ClientData, drMsg string) error {
	queueName := fmt.Sprintf("%s-dr-msg", clientData.ClientName)
	contentType := "text/plain"
	err := drRep.conn.Send(queueName, contentType, []byte(drMsg), stomp.SendOpt.Header("destication-type", "ANYCAST"))
	if err != nil {
		return err
	}
	return nil
}
