package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"github.com/MooooonStar/mixin-sdk-go/messenger"
	mixin "github.com/MooooonStar/mixin-sdk-go/network"
)

type Listener struct {
	*messenger.Messenger
}

// interface to implement if you want to handle the message
func (l *Listener) OnMessage(ctx context.Context, msg messenger.MessageView, userId string) error {
	data, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		return err
	}
	if msg.Category == "SYSTEM_ACCOUNT_SNAPSHOT" {
		var transfer messenger.TransferView
		if err := json.Unmarshal(data, &transfer); err != nil {
			return err
		}
		log.Println("I got a coin: ", transfer.Amount)
		mixin.Transfer(msg.UserId,transfer.Amount,transfer.AssetId,"",messenger.UuidNewV4().String(),
									PinCode,PinToken,UserId,SessionId,PrivateKey)
		return nil
		// return l.SendPlainText(ctx, msg.ConversationId, msg.UserId, string(data))
	} else if msg.Category == "PLAIN_TEXT" {
		log.Printf("I got a message, it said: %s", string(data))
		return l.SendPlainText(ctx, msg.ConversationId, msg.UserId, string(data))
	} else {
		log.Println("Unknown message!", msg.Category)
		return err
	}
}
const (
	UserId    = "21042518-85c7-4903-bb19-f311813d1f51"
	PinCode   = "911424"
	SessionId = "4267b63d-3daa-449e-bc13-970aa0357776"
	PinToken  = "gUUxpm3fPRVkKZNwA/gk10SHHDtR8LmxO+N6KbsZ/jymmwwVitUHKgLbk1NISdN8jBvsYJgF/5hbkxNnCJER5XAZ0Y35gsAxBOgcFN8otsV6F0FAm5TnWN8YYCqeFnXYJnqmI30IXJTAgMhliLj7iZsvyY/3htaHUUuN5pQ5F5s="
	//please delele the blank of PrivateKey the before each line
	PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQCDXiWJRLe9BzPtXmcVe6acaFTY9Ogb4Hc2VHFjKFsp7QRVCytx
3KC/LRojTFViwwExaANTZQ6ectwpAxIvzeYeHDZCXCh6JRFIYK/ZuREmYPcPQEWD
s92Tv/4XTAdTH8l9UJ4VQY4zwqYMak237N9xEvowT0eR8lpeJG0jAjN97QIDAQAB
AoGADvORLB1hGCeQtmxvKRfIr7aEKak+HaYfi1RzD0kRjyUFwDQkPrJQrVGRzwCq
GzJ8mUXwUvaGgmwqOJS75ir2DL8KPz7UfgQnSsHDUwKqUzULgW6nd/3OdDTYWWaN
cDjbkEpsVchOpcdkywvZhhyGXszpM20Vr8emlBcFUOTfpTUCQQDVVjkeMcpRsImV
U3tPYyiuqADhBTcgPBb+Ownk/87jyKF1CZOPvJAebNmpfJP0RMxUVvT4B9/U/yxZ
WNLhLtCXAkEAnaOEuefUxGdE8/55dUTEb7xrr22mNqykJaax3zFK+hSFBrM3gUY5
fEETtHnl4gEdX4jCPybRVc1JSFY/GWoyGwJBAKoLti95JHkErEXYavuWYEEHLNwv
mgcZnoI6cOKVfEVYEEoHvhTeCkoWHVDZOd2EURIQ1eY18JYIZ0M4Z66R8DUCQCsK
iKTR3dA6eiM8qiEQw6nWgniFscpf3PnCx/Iu3U/m5mNr743GhM+eXSj7136b209I
YfEoQiPxRz8O/W+NBV0CQQDVPNqJlFD34MC9aQN42l3NV1hDsl1+nSkWkXSyhhNR
MpobtV1a7IgJGyt5HxBzgNlBNOayICRf0rRjvCdw6aTP
-----END RSA PRIVATE KEY-----`
)


func main() {
	ctx := context.Background()
	m := messenger.NewMessenger(UserId, SessionId, PrivateKey)
	//replace with your own listener
	//go m.Run(ctx, messenger.DefaultBlazeListener{})
	l := &Listener{m}
	go m.Run(ctx, l)

	// //your mixin user id, can get from  "Search User"
	// snow := "7b3f0a95-3ee9-4c1b-8ae9-170e3877d909"
	//
	// //must create conversation first. If have created before, skip this step.
	// if _, err := m.CreateConversation(ctx, messenger.CategoryContact, messenger.Participant{UserID: snow}); err != nil {
	// 	log.Println("create conversation error", err)
	// }
	//
	// conversation, err := m.CreateConversation(ctx, messenger.CategoryContact,
	// 	messenger.Participant{UserID: snow},
	// )
	// if err != nil {
	// 	log.Println("create error", err)
	// }
	//
	// if err := m.SendImage(ctx, conversation.ID, snow, "../donate.png"); err != nil {
	// 	log.Println("send image error:", err)
	// }
	//
	// if err := m.SendVideo(ctx, conversation.ID, snow, "../sample.mp4"); err != nil {
	// 	log.Println("send video error", err)
	// }
	//
	// if err := m.SendFile(ctx, conversation.ID, snow, "../demo.pdf", "application/pdf"); err != nil {
	// 	log.Println("send video error", err)
	// }
	//
	// if err := m.SendPlainText(ctx, conversation.ID, snow, "please send me a message and transfer some CNB to me."); err != nil {
	// 	log.Println("send text error:", err)
	// }

	select {}
}
