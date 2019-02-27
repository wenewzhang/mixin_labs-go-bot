package main
import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"
	NewBlazeClient "github.com/MooooonStar/mixin-sdk-go/messenger"
  config "github.com/wenewzhang/mixin_labs-go-bot/config"
	mixin "github.com/MooooonStar/mixin-sdk-go/network"
)
type BlazeClient =  messenger.BlazeClient
type BlazeListener = messenger.BlazeListener
// type NewBlazeClient = messenger.NewBlazeClient
// Messenger mixin messenger
type Messenger struct {
	*mixin.User
	*BlazeClient
}

// NewMessenger create messenger
func NewMessenger(userId, sessionId, privateKey string) *Messenger {
	user := mixin.NewUser(userId, sessionId, privateKey)
	client := NewBlazeClient(userId, sessionId, privateKey)
	return &Messenger{user, client}
}

func (m *Messenger) Run(ctx context.Context, listener BlazeListener) {
	for {
		if err := m.Loop(ctx, listener); err != nil {
			log.Println("Blaze server error", err)
			time.Sleep(1 * time.Second)
		}
		m.BlazeClient = NewBlazeClient(m.UserId, m.SessionId, m.PrivateKey)
	}
}

type DefaultBlazeListener struct{}

// interface to implement if you want to handle the message
func (l DefaultBlazeListener) OnMessage(ctx context.Context, msg MessageView, userId string) error {
	data, err := base64.StdEncoding.DecodeString(msg.Data)
	// log.Printf("on message",msg.action)
	log.Println(msg)
	if err != nil {
		return err
	}
	if msg.Category == "SYSTEM_ACCOUNT_SNAPSHOT" {
		var transfer TransferView
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

func main() {
	ctx := context.Background()
	m := NewMessenger(config.UserId, config.SessionId, config.PrivateKey)
}
