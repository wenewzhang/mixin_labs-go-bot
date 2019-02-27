package main

import (
	"context"
  config "github.com/wenewzhang/mixin_labs-go-bot/config"
  msgbot "github.com/wenewzhang/mixin_labs-go-bot/msgbot"
	// "github.com/MooooonStar/mixin-sdk-go/messenger"
)



func main() {
	ctx := context.Background()
	m := msgbot.NewMessenger(config.UserId, config.SessionId, config.PrivateKey)
	//replace with your own listener
	go m.Run(ctx, msgbot.DefaultBlazeListener{})

	select {}
}
