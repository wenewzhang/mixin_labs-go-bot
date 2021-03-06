# Go Bitcoin tutorial based on Mixin Network : Receive and send Bitcoin in Mixin Messenger
![cover](https://github.com/wenewzhang/mixin_labs-go-bot/raw/master/Bitcoin_go.jpg)
In [the previous chapter](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README.md), we created our first app, when user sends "Hello,world!", the bot reply the same message.

> main.go

```go
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
		if string(data) == "g" {
			payLinkEOS := "https://mixin.one/pay?recipient=" +
							 msg.UserId  + "&asset=" +
							 "6cfe566e-4aad-470b-8c9a-2fd35b49c68d"   +
							 "&amount=" + "0.1" +
							 "&trace="  + messenger.UuidNewV4().String() +
							 "&memo="
		  payLinkBTC := "https://mixin.one/pay?recipient=" +
							 msg.UserId  + "&asset=" +
							 "c6d0c728-2624-429b-8e0d-d9d19b6592fa"   +
							 "&amount=" + "0.001" +
							 "&trace="  + messenger.UuidNewV4().String() +
							 "&memo="
		  log.Println(payLinkBTC)
			BtnEOS := messenger.Button{Label: "Pay EOS 0.1", Color: "#0080FF", Action: payLinkEOS}
			BtnBTC := messenger.Button{Label: "Pay BTC 0.001", Color: "#00FF80", Action: payLinkBTC}
			if err := l.SendAppButtons(ctx, msg.ConversationId, msg.UserId, BtnEOS, BtnBTC); err != nil {
				return err
			}
			return nil
		} else if string(data) == "a"  {
			card := messenger.AppCard{Title: "CNB", Description: "Chui Niu Bi", Action: "http://www.google.cn",
				IconUrl: "https://images.mixin.one/0sQY63dDMkWTURkJVjowWY6Le4ICjAFuu3ANVyZA4uI3UdkbuOT5fjJUT82ArNYmZvVcxDXyNjxoOv0TAYbQTNKS=s128"}
			if err := l.SendAppCard(ctx, msg.ConversationId, msg.UserId, card); err != nil {
				return err
			}
			return nil
		} else if string(data) == "r" {
			mixin.Transfer(msg.UserId,"0.0001","c6d0c728-2624-429b-8e0d-d9d19b6592fa","",messenger.UuidNewV4().String(),
										PinCode,PinToken,UserId,SessionId,PrivateKey)
			return nil
		} else { return l.SendPlainText(ctx, msg.ConversationId, msg.UserId, string(data)) }
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
	select {}
}


```
### Hello Bitcoin!
Build the bot and execute it in the project directory.
```bash
cd mixin_labs-go-bot
go build
./mixin_labs-go-bot
```

Developer can send Bitcoin to their bots in message panel. The bot receive the Bitcoin and then send back immediately.
![transfer and tokens](https://github.com/wenewzhang/mixin_network-nodejs-bot2/blob/master/transfer-any-tokens.jpg)

User can pay 0.001 Bitcoin to bot by click the button and the 0.001 Bitcoin will be refunded in 1 second,In fact, user can pay any coin either.
![pay-link](https://github.com/wenewzhang/mixin_network-nodejs-bot2/blob/master/Pay_and_refund_quickly.jpg)

## Source code summary
```go
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
}
```
Call mixin.Transfer to refund the coins back to user.

## Advanced usage

#### APP_BUTTON_GROUP
In some payment scenario, for example:
The coin exchange provides coin-exchange service which transfer BTC to EOS ETH, BCH etc,
you want show the clients many pay links with different amount, APP_BUTTON_GROUP can help you here.
```go
payLinkEOS := "https://mixin.one/pay?recipient=" +
         msg.UserId  + "&asset=" +
         "6cfe566e-4aad-470b-8c9a-2fd35b49c68d"   +
         "&amount=" + "0.1" +
         "&trace="  + messenger.UuidNewV4().String() +
         "&memo="
payLinkBTC := "https://mixin.one/pay?recipient=" +
         msg.UserId  + "&asset=" +
         "c6d0c728-2624-429b-8e0d-d9d19b6592fa"   +
         "&amount=" + "0.001" +
         "&trace="  + messenger.UuidNewV4().String() +
         "&memo="
log.Println(payLinkBTC)
BtnEOS := messenger.Button{Label: "Pay EOS 0.1", Color: "#0080FF", Action: payLinkEOS}
BtnBTC := messenger.Button{Label: "Pay BTC 0.001", Color: "#00FF80", Action: payLinkBTC}
if err := l.SendAppButtons(ctx, msg.ConversationId, msg.UserId, BtnEOS, BtnBTC); err != nil {
  return err
}
```
Here show clients two buttons for EOS and BTC, you can add more buttons in this way.

#### APP_CARD
Maybe a group of buttons too simple for you, try a pay link which show a icon: APP_CARD.
```go
card := messenger.AppCard{Title: "CNB", Description: "Chui Niu Bi", Action: "http://www.google.cn",
  IconUrl: "https://images.mixin.one/0sQY63dDMkWTURkJVjowWY6Le4ICjAFuu3ANVyZA4uI3UdkbuOT5fjJUT82ArNYmZvVcxDXyNjxoOv0TAYbQTNKS=s128"}
if err := l.SendAppCard(ctx, msg.ConversationId, msg.UserId, card); err != nil {
  return err
}
```
![APP_CARD](https://github.com/wenewzhang/mixin_labs-python-bot/raw/master/app_card.jpg)

[Full source code](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/main.go)
