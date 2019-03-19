package main

import (
    "encoding/base64"
    "encoding/json"
    "github.com/satori/go.uuid"
    "github.com/vmihailenco/msgpack"
    "fmt"
    "encoding/csv"
    "log"
    "os"
    "bufio"
    mixin "github.com/MooooonStar/mixin-sdk-go/network"
)

const (
	UserId    = "21042518-85c7-4903-bb19-f311813d1f51"
	PinCode   = "540806"
	SessionId = "806a2ea0-e1e0-409f-beac-d6410af7bf93"
	PinToken  = "N3HqYZzHQ67/bMXjm8M9NotKsxWabSWEy1mYICFiaFPUQaRkWa+AiUQsEoKoeRtnbH4sy912vn9/RII//yO9GF94MAACWSShUubbADOmJjGKvgxBRSMHi9h7I4jCh1LQhkUq7iUKVY5GBGCWp83ZbcyvEGWlBnvOcqb010jhMpE="
	//please delele the blank of PrivateKey the before each line
	PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQCOJB1fIMpyPCh4pPM0cZMhYfDAA2IwgYrjwYFT0EGlkjJW8TK5
BIIQXvE83VHUSmuVFPOgpOyycbLvxwfau2TP6FTLLN7WEgyVyA4CzpiQ0ihPHXjQ
TtUJKnJ7Uj3dabG+0BnkvksKSpnhXv+lBCEd7Y8P+f2IcXAURECSbKIrdwIDAQAB
AoGBAIFo6gDY5tgYYzRr4SznFnA3LixzKCtHVb9ERs2a9pmawBAd6vM94nirJ/El
AsJHuyjw+VpRrVpNX/8j8se3AvFUPF4Q67sMOAob3dsOHy3SKIysOKwRzLuYF69m
rweJPSXJuqFWc6pHuXHArW/eL5ZlWbh1dKdU+2EYEAy/kcjpAkEA2zW8qt6zv5D5
IjvCIz3mtincdng9DuZE5SzsP7fxRkEMVqP4QLrRnfwuHYuZj9KjuEFX85M2iKot
cv8KHI0ovQJBAKX/JadZmoL4HJKz/OfXH3lU6zy787gY0xYcDPh4MSMDNVSOcqhe
lOnintQD5vPaGZUl6aBkO5Dmp6JIgELfqkMCQQDRdbGHnDEpVT+ZJHzG6/kuCyXr
1cySFhmy2pAL+pmDRdiiWR93yotNaJAwDxp2wRFLmLSPvBUZ1XKENYrV6VQJAkBV
adY8KDUDExvQuOB3gw/k5LcuRx//KHblN4XNDDtsYqg8XBfPXuuM9Vj4ixF5hE4J
mrp+F1U3GBhFvryQrHn1AkEAlDLpBgaxt4kpCeXUr/9lUoKsF/8taOlS3IZjRbBJ
A3BlaWdHIvUHhqpVbfeCYv7m3GnIs+Kfo1I56haIRVFrNw==
-----END RSA PRIVATE KEY-----`
)

type OrderAction struct {
    A  uuid.UUID  // asset uuid
}
func ReadAssetInfo(asset_id string) ( map[string]interface{}, string) {
  var UserInfoMap map[string]interface{}
  csvFile, err := os.Open("mybitcoin_wallet.csv")
  if err != nil {
         log.Fatal(err)
  }
  reader := csv.NewReader(bufio.NewReader(csvFile))
  record, err := reader.Read()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(record[3])
  PrivateKey2           := record[0]
  SessionID2      		  := record[2]
  UserID2               := record[3]
  UserInfoBytes, err    := mixin.ReadAsset(mixin.GetAssetId(asset_id),
                                         UserID2,SessionID2,PrivateKey2)
  if err != nil {
          log.Fatal(err)
  }
  // fmt.Println(string(UserInfoBytes))
  if err := json.Unmarshal(UserInfoBytes, &UserInfoMap); err != nil {
      panic(err)
  }
  csvFile.Close()
  return UserInfoMap, UserID2
}

func main() {
  // Pack memo
  packUuid, _ := uuid.FromString("c6d0c728-2624-429b-8e0d-d9d19b6592fa")
  pack, _ := msgpack.Marshal(OrderAction{A: packUuid,})
  memo := base64.StdEncoding.EncodeToString(pack)
  // gaFBxBDG0McoJiRCm44N2dGbZZL6
  fmt.Println(memo)
  // Parse memo
  parsedpack, _ := base64.StdEncoding.DecodeString(memo)
  orderAction := OrderAction{}
  _ = msgpack.Unmarshal(parsedpack, &orderAction)
  fmt.Println(orderAction.A)

  scanner   := bufio.NewScanner(os.Stdin)
	var PromptMsg string
	PromptMsg  = "1: Create Wallet\n2: Read Bitcoin balance & Address \n3: Read USDT balance & Address\n4: Read EOS balance & address\n"
	PromptMsg += "5: Read EOS address\n6: Transfer Bitcoin from bot to new account\n7: Transfer Bitcoin from new account to Master\n"
	PromptMsg += "8: Withdraw bot's Bitcoin\na: Verify Pin\nd: Create Address and Delete it\nr: Create Address and read it\n"
	PromptMsg += "q: Exit \nMake your choose:"
	for  {
	 fmt.Print(PromptMsg)
	 var cmd string
	 scanner.Scan()
	 cmd = scanner.Text()
	 if cmd == "q" {
			 break
	 }
  if cmd == "1" {
    fo, err := os.OpenFile("mybitcoin_wallet.csv",
                           os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
      panic(err)
      return
    }

    user,err := mixin.CreateAppUser("Tom cat", PinCode, UserId,
                                   SessionId, PrivateKey)
    if err != nil {
        panic(err)
    }
    records := [][]string {
                        {user.PrivateKey,
                          user.PinToken,
                          user.SessionId,
                          user.UserId,
                          user.PinCode},
                        }
    w := csv.NewWriter(fo)
    if err := w.Error(); err != nil {
      log.Fatalln("error writing csv:", err)
    }
    w.WriteAll(records) // calls Flush internally
    fo.Close()
    log.Println(user)
  }
  if cmd == "2" {
    userInfo, userID := ReadAssetInfo("BTC")
    fmt.Println("User ID ",userID, "'s BTC Address is: ",
               userInfo["data"].(map[string]interface{})["public_key"])
    fmt.Println("Balance is: ",
               userInfo["data"].(map[string]interface{})["balance"])
  }
  if cmd == "3" {
    userInfo, userID := ReadAssetInfo("USDT")
    fmt.Println("User ID ",userID, "'s USDT Address is: ",
               userInfo["data"].(map[string]interface{})["public_key"])
    fmt.Println("Balance is: ",
               userInfo["data"].(map[string]interface{})["balance"])
  }
  if cmd == "4" {
    userInfo, userID := ReadAssetInfo("EOS")
    fmt.Println("User ID ",userID, "'s EOS Address is: ",
               userInfo["data"].(map[string]interface{})["account_name"],
               userInfo["data"].(map[string]interface{})["account_tag"])
    fmt.Println("Balance is: ",
               userInfo["data"].(map[string]interface{})["balance"])
  }
  }
  // c6d0c728-2624-429b-8e0d-d9d19b6592fa
}
