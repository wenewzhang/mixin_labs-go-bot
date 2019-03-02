# 基于Mixin Network的Go语言比特币开发教程:创建机器人
![cover](https://github.com/wenewzhang/mixin_labs-go-bot/raw/master/Bitcoin_go.jpg)
[Mixin Network](https://mixin.one) 是一个免费的 极速的端对端加密数字货币交易系统.
在本章中，你可以按教程在Mixin Messenger中创建一个bot来接收用户消息, 学到如何给机器人转**比特币** 或者 让机器人给你转**比特币**.

[Mixin Network的开发资源汇编](https://github.com/awesome-mixin-network/index_of_Mixin_Network_resource)

## 课程简介
1. [创建一个接受消息的机器人](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README-zhchs.md)
2. [机器人接受比特币并立即退还用户](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README2-zhchs.md)
3. [创建比特币钱包](#)

## 创建一个接受消息的机器人

通过本教程，你将学会如何用Go创建一个机器人APP,让它能接受消息.

### Go 1.12 的安装:
从Go官网下载安装 [Go](https://golang.org/dl/)

macOS

下载安装包 [go1.12.darwin-amd64.pkg](https://dl.google.com/go/go1.12.darwin-amd64.pkg) 双击运行，然后按提示安装, 最后将Go的bin目录加入到$PATH中.
```bash
echo 'export PATH="/usr/local/opt/go/libexec/bin:$PATH"' >> ~/.bash_profile
source  ~/.bashrc
```
如果一切正常，运行 **go version**就可以看到如下提示了！
```bash
go version
go version go1.11.5 darwin/amd64
```

Ubuntu 18.04
```bash
root@n3:/usr/local/bin# snap install go --classic
```

Ubuntu 16.04
```bash
mkdir /usr/local/src
wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
tar xvf go1.12.linux-amd64.tar.gz
echo 'export PATH=/usr/local/src/go/bin:$PATH' >> ~/.bashrc
root@n3:/usr/local/src# source  ~/.bashrc
```

如果一切正常，运行 **go version**就可以看到如下提示了！
```bash
root@n3:/usr/local/bin# go version
go version go1.12 linux/amd64
```

### 创建Go的工作目录
强烈推荐为Go创建一个工作目录，这让你少了很多关于包的引用的麻烦。
macOS
```bash
mkdir ~/workspace/go
echo 'export GOPATH="$HOME/workspace/go"' >> ~/.bash_profile
source ~/.bash_profile
```

Ubuntu
```bash
mkdir ~/workspace/go
echo 'export GOPATH="$HOME/workspace/go"' >> ~/.bashrc
source ~/.bash_profile
```

### 安装 Mixin Network SDK for Go
```bash
root@n3:~# go get github.com/MooooonStar/mixin-sdk-go
package github.com/MooooonStar/mixin-sdk-go: no Go files in /root/workspace/go/src/github.com/MooooonStar/mixin-sdk-go
```

不用担心 "no Go files" 的提示, **ls** 可以找到如下目录与文件，事实上，SDK分别在messenger,network中.
```bash
ls $GOPATH/src/github.com/MooooonStar/mixin-sdk-go

README.md	messenger	network
```

### 在GOPATH下创建项目目录
```bash
cd ~/workspace/go/src
mkdir mixin_labs-go-bot
cd mixin_labs-go-bot

```

## 你好，世界!

### 创建第一个机器人APP
按下面的提示，到[mixin.one](https://mixin.one)创建一个APP[tutorial](https://mixin-network.gitbook.io/mixin-network/mixin-messenger-app/create-bot-account).

### 生成相应的参数
记下这些[生成的参数](https://mixin-network.gitbook.io/mixin-network/mixin-messenger-app/create-bot-account#generate-secure-parameter-for-your-app)
它们将用于main.go中.

![mixin_network-keys](https://github.com/wenewzhang/mixin_labs-php-bot/raw/master/mixin_network-keys.jpg)
在项目目录下，创建main.go,将生成的参数，替换成你的！
> main.go
```go
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
```
将上面的参数，替换成你在mixin.one生成的。

完整而又简洁的代码如下
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
  if msg.Category == "PLAIN_TEXT" {
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
	l := &Listener{m}
	go m.Run(ctx, l)
	select {}
}
```
### 编译与运行
执行 **go build** 将创建一个mixin_labs-go-bot,然后执行
```bash
cd mixin_labs-go-bot
go build
./mixin_labs-go-bot
```
在手机安装 [Mixin Messenger](https://mixin.one/),增加机器人为好友,(比如这个机器人是7000101639) 然后发送消息给它,效果如下!

![mixin_messenger](https://github.com/wenewzhang/mixin_labs-php-bot/raw/master/helloworld.jpeg)

## 源代码解释
WebSocket是建立在TCP基础之上的全双工通讯方式，连接到Mixin Network并发送"LISTPENDINGMESSAGES"消息，服务器以后会将收到的消息转发给此程序!
```go
ctx := context.Background()
m := messenger.NewMessenger(UserId, SessionId, PrivateKey)
l := &Listener{m}
go m.Run(ctx, l)
```
当服务器给机器人推送消息的时候,机器人会原封不动的回复回去.
```go
func (l *Listener) OnMessage(ctx context.Context, msg messenger.MessageView, userId string) error {
	data, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		return err
	}
  if msg.Category == "PLAIN_TEXT" {
		log.Printf("I got a message, it said: %s", string(data))
		return l.SendPlainText(ctx, msg.ConversationId, msg.UserId, string(data))
	} else {
		log.Println("Unknown message!", msg.Category)
		return err
	}
}
```
Mixin Messenger支持的消息类型很多，除了文本，还有图片，视频，语音等等，具体可到下面链接查看:  [WebSocket消息类型](https://developers.mixin.one/api/beta-mixin-message/websocket-messages/).

### 完成
现在你的机器人APP运行起来了，你打算如何改造你的机器人呢？

完整的代码[在这儿](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/main.go)
