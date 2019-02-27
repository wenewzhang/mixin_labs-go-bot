package main

import (
	"log"
  config "github.com/wenewzhang/mixin_labs-go-bot/config"
	mixin "github.com/MooooonStar/mixin-sdk-go/network"
)

func main() {
	user := mixin.NewUser(config.UserId, config.SessionId, config.PrivateKey, config.PinCode, config.PinToken)
	//user := mixin.NewUser(UserId, SessionId, PrivateKey)

	profile, err := user.ReadProfile()
	if err != nil {
		log.Fatal("Read profile error", err)
	}
	log.Println("profile", string(profile))

	assets, err := user.ReadAssets()
	if err != nil {
		log.Fatal("Read assets error", assets)
	}
	log.Println("assets", string(assets))
}
