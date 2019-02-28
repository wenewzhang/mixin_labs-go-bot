# Go Mixin Messenger application development tutorial
![cover](#)
A Mixin messenger bot will be created in this tutorial. The bot is powered by Go, it echo message and Bitcoin from user.

Full Mixin network resource [index](https://github.com/awesome-mixin-network/index_of_Mixin_Network_resource)

## What you will learn from this tutorial
1. [How to create bot in Mixin messenger and reply message to user](https://github.com/wenewzhang/mixin_labs-go-bot#create-bot-and-receive-message-from-user)
2. [How to receive Bitcoin and send Bitcoin in Mixin Messenger](#)
3. [How to create a Bitcoin wallet based on Mixin Network API](#)
## How to create bot in Mixin messenger and reply message to user
## Go 1.12 installation:
Download the latest [Go](https://golang.org/dl/) from here.

macOS

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

Check Go's installation
```bash
root@n3:/usr/local/bin# go version
go version go1.12 linux/amd64
```
or
```bash
go version
go version go1.11.5 darwin/amd64
```

## Set Go 's workspace
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

## Install the Go SDK
```bash
root@n3:~# go get github.com/MooooonStar/mixin-sdk-go
package github.com/MooooonStar/mixin-sdk-go: no Go files in /root/workspace/go/src/github.com/MooooonStar/mixin-sdk-go
```

## Make project of this bot
```bash
cd ~/workspace/go/src
mkdir mixin_labs-go-bot
cd mixin_labs-go-bot

```

## Hello, world in Go

### Create your first app in Mixin Network developer dashboard
You need to create an app in dashboard. This [tutorial](https://mixin-network.gitbook.io/mixin-network/mixin-messenger-app/create-bot-account) can help you.

### Generate parameter of your app in dashboard
After app is created in dashboard, you still need to [generate parameter](https://mixin-network.gitbook.io/mixin-network/mixin-messenger-app/create-bot-account#generate-secure-parameter-for-your-app)
and write down required content, these content will be written into main.go file.

![mixin_network-keys](https://github.com/wenewzhang/mixin_labs-php-bot/blob/master/mixin_network-keys.jpg)
In the folder, create a file: main.go  Copy the following content into it.
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
Replace the value with content generated in dashboard.
