package main

import (
	"context"
	"log"
  config "github.com/wenewzhang/mixin_labs-go-bot/config"
	"github.com/MooooonStar/mixin-sdk-go/messenger"
)

func main() {
	ctx := context.Background()
	m := messenger.NewMessenger(config.UserId, config.SessionId, config.PrivateKey)
	//replace with your own listener
	go m.Run(ctx, messenger.DefaultBlazeListener{})

	snow := "0b4f49dc-8fb4-4539-9a89-fb3afc613747"

	//must create conversation first. If have created before, skip this step.
	if _, err := m.CreateConversation(ctx, messenger.CategoryContact, messenger.Participant{UserID: snow}); err != nil {
		log.Println("create conversation error", err)
	}
	conversation, err := m.CreateConversation(ctx, messenger.CategoryGroup,
		messenger.Participant{UserID: snow},
	)
	if err != nil {
		log.Println("create error", err)
	}

	if err := m.SendPlainText(ctx, conversation.ID, snow, "please send me a message."); err != nil {
		log.Println("send text error:", err)
	}

	if err := m.SendImage(ctx, conversation.ID, snow, "../donate.png"); err != nil {
		log.Println("send image error:", err)
	}

	if err := m.SendVideo(ctx, conversation.ID, snow, "../sample.mp4"); err != nil {
		log.Println("send video error", err)
	}

	if err := m.SendFile(ctx, conversation.ID, snow, "../demo.pdf", "application/pdf"); err != nil {
		log.Println("send video error", err)
	}

	select {}
}
