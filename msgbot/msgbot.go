package msgbot

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"
  "github.com/MooooonStar/mixin-sdk-go/messenger"
	mixin "github.com/MooooonStar/mixin-sdk-go/network"
)

// Messenger mixin messenger
type Messenger struct {
	*mixin.User
	*messenger.BlazeClient
}

// NewMessenger create messenger
func NewMessenger(userId, sessionId, privateKey string) *Messenger {
	user := mixin.NewUser(userId, sessionId, privateKey)
	client := messenger.NewBlazeClient(userId, sessionId, privateKey)
	return &Messenger{user, client}
}

func (m *Messenger) Run(ctx context.Context, listener messenger.BlazeListener) {
	for {
		if err := m.Loop(ctx, listener); err != nil {
			log.Println("Blaze server error", err)
			time.Sleep(1 * time.Second)
		}
		m.BlazeClient = messenger.NewBlazeClient(m.UserId, m.SessionId, m.PrivateKey)
	}
}

type DefaultBlazeListener struct{}

// interface to implement if you want to handle the message
func (l DefaultBlazeListener) OnMessage(ctx context.Context, msg messenger.MessageView, userId string) error {
	data, err := base64.StdEncoding.DecodeString(msg.Data)
	log.Println(msg)
	if err != nil {
		return err
	}
	if msg.Category == "SYSTEM_ACCOUNT_SNAPSHOT" {
		var transfer messenger.TransferView
		if err := json.Unmarshal(data, &transfer); err != nil {
			return err
		}
		log.Println("I got a snapshot: ", transfer)
		return nil
	} else {
		log.Printf("I got a message, it said: %s", string(data))
		return nil
	}
}
