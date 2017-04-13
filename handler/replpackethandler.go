package replpackethandler

import (
	"fmt"
	"log"

	"github.com/piot/hasty-protocol/channel"
	"github.com/piot/hasty-protocol/commands"
)

type Handler struct {
}

func (in *Handler) HandleConnect(cmd commands.Connect) error {
	log.Printf("repl:%s", cmd)
	return fmt.Errorf("Repl doesn't support connect")
}

func (in *Handler) HandlePublishStream(cmd commands.PublishStream) error {
	log.Printf("repl:%s", cmd)
	return fmt.Errorf("Repl doesn't support publish")
}

func (in *Handler) HandleSubscribeStream(cmd commands.SubscribeStream) {
	log.Printf("repl:%s", cmd)
}

func (in *Handler) HandlePing(cmd commands.Ping) {
	log.Printf("repl:%s", cmd)
}

func (in *Handler) HandlePong(cmd commands.Pong) {
	log.Printf("repl:%s", cmd)
}

func (in *Handler) HandleLogin(cmd commands.Login) error {
	log.Printf("repl:%s", cmd)
	return nil
}

func (in *Handler) HandleUnsubscribeStream(cmd commands.UnsubscribeStream) {
	log.Printf("repl:%s", cmd)

}

func (in *Handler) HandleCreateStream(cmd commands.CreateStream) (channel.ID, error) {
	log.Printf("repl:%s", cmd)

	return channel.ID{}, fmt.Errorf("Repl can't handle create stream")
}

func (in *Handler) HandleStreamData(cmd commands.StreamData) {
	log.Printf("repl:%s", cmd)
}

func (in *Handler) HandleTransportDisconnect() {
	log.Printf("Handle Transport Disconnect")
}
