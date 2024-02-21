package provider

import (
	"fmt"

	"transmitter-artemis/config"

	"github.com/go-stomp/stomp/v3"
)

// NewArtemis create a new stomp connection
func NewArtemis() (conn *stomp.Conn, err error) {
	cfg := config.Configuration.Artemis

	conn, err = stomp.Dial(
		"tcp",
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		stomp.ConnOpt.Login(cfg.Username, cfg.Password),
	)

	return
}
