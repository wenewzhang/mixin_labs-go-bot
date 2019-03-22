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
		"io/ioutil"
		"net/http"
		"bytes"
		"time"
    mixin "github.com/MooooonStar/mixin-sdk-go/network"
		"github.com/MooooonStar/mixin-sdk-go/messenger"
)

const (
	UserId     = "21042518-85c7-4903-bb19-f311813d1f51"
	PinCode    = "540806"
	SessionId  = "806a2ea0-e1e0-409f-beac-d6410af7bf93"
	PinToken   = "N3HqYZzHQ67/bMXjm8M9NotKsxWabSWEy1mYICFiaFPUQaRkWa+AiUQsEoKoeRtnbH4sy912vn9/RII//yO9GF94MAACWSShUubbADOmJjGKvgxBRSMHi9h7I4jCh1LQhkUq7iUKVY5GBGCWp83ZbcyvEGWlBnvOcqb010jhMpE="
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
	EXIN_BOT   = "61103d28-3ac2-44a2-ae34-bd956070dab1"
	// EXIN_BOT   = "0b4f49dc-8fb4-4539-9a89-fb3afc613747"
)

type OrderAction struct {
    A  uuid.UUID  // asset uuid
}

type Error struct {
	Status      int    `json:"status"`
	Code        int    `json:"code"`
	Description string `json:"description"`
	trace       string
}

// type respData struct {
// 	Zero      string `json:"0"`
// 	One string `json:"1"`
// 	Two       string `json:"2"`
// }

func (e Error) Error() string {
	bt, _ := json.Marshal(e)
	return string(bt)
}

var httpClient = &http.Client{Timeout: time.Duration(10 * time.Second)}

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
  fmt.Println(string(UserInfoBytes))
  if err := json.Unmarshal(UserInfoBytes, &UserInfoMap); err != nil {
      panic(err)
  }
  csvFile.Close()
  return UserInfoMap, UserID2
}

func GetWalletInfo() ( string, string, string, string, string) {
  csvFile, err := os.Open("mybitcoin_wallet.csv")
  if err != nil {
         log.Fatal(err)
  }
  reader := csv.NewReader(bufio.NewReader(csvFile))
  record, err := reader.Read()
  if err != nil {
    log.Fatal(err)
  }
  csvFile.Close()
  return record[0], record[1], record[2], record[3], record[4]
}

func GetMarketPrice(asset_id string) ([]byte, error)  {
	var body []byte
	req, err := http.NewRequest("GET", "https://exinone.com/exincore/markets?base_asset="+asset_id, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
  // fmt.Println(resp.Body)
	bt, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		var resp struct {
			Error Error `json:"error"`
		}
		err = json.Unmarshal(bt, &resp)
		if err == nil {
			err = resp.Error
		}
	}
	return bt, err
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
	PromptMsg += "5: pay 0.0001 BTC buy USDT\n6: Read ExinCore Price(USDT)\n7: Read ExinCore Price(BTC)\n"
	PromptMsg += "8: pay 1 USDT buy BTC\n9: Read Snapshots\na: Verify bot PIN code\nv: Verify wallet PIN code\n"
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
	if cmd == "5" {
		packUuid, _ := uuid.FromString(mixin.GetAssetId("USDT"))
		pack, _ := msgpack.Marshal(OrderAction{A: packUuid})
		memo := base64.StdEncoding.EncodeToString(pack)
		// fmt.Println(memo)
		priKey, pToken, sID, userID, uPIN := GetWalletInfo()
		bt, err := mixin.Transfer(EXIN_BOT,"0.0001",mixin.GetAssetId("BTC"),memo,
														 messenger.UuidNewV4().String(),
														 uPIN,pToken,userID,sID,priKey)
		if err != nil {
						log.Fatal(err)
		}
		fmt.Println(string(bt))
  }
  if cmd == "6" {
		priceInfo, err := GetMarketPrice(mixin.GetAssetId("USDT"))
		if err != nil {
			log.Fatal(err)
		}

		var marketInfo map[string]interface{}
		err = json.Unmarshal([]byte(priceInfo), &marketInfo)
    fmt.Println("Asset | price | min | max | exchanges")
		for _, v := range (marketInfo["data"].(map[string]interface{})) {
			fmt.Println(v.(map[string]interface{})["exchange_asset_symbol"],"/",
									v.(map[string]interface{})["base_asset_symbol"],
									v.(map[string]interface{})["price"],
									v.(map[string]interface{})["minimum_amount"],
									v.(map[string]interface{})["maximum_amount"],
									v.(map[string]interface{})["exchanges"],
								 )
		}
	}
	if cmd == "7" {
		priceInfo, err := GetMarketPrice(mixin.GetAssetId("BTC"))
		if err != nil {
			log.Fatal(err)
		}

		var marketInfo map[string]interface{}
		err = json.Unmarshal([]byte(priceInfo), &marketInfo)
    fmt.Println("Asset | price | min | max | exchanges")
		for _, v := range (marketInfo["data"].(map[string]interface{})) {
			fmt.Println(v.(map[string]interface{})["exchange_asset_symbol"],"/",
									v.(map[string]interface{})["base_asset_symbol"],
									v.(map[string]interface{})["price"],
									v.(map[string]interface{})["minimum_amount"],
									v.(map[string]interface{})["maximum_amount"],
									v.(map[string]interface{})["exchanges"],
								 )
		}
	}
	if cmd == "8" {
		packUuid, _ := uuid.FromString(mixin.GetAssetId("BTC"))
		pack, _ := msgpack.Marshal(OrderAction{A: packUuid})
		memo := base64.StdEncoding.EncodeToString(pack)
		// fmt.Println(memo)
		priKey, pToken, sID, userID, uPIN := GetWalletInfo()
		bt, err := mixin.Transfer(EXIN_BOT,"1",mixin.GetAssetId("USDT"),memo,
														 messenger.UuidNewV4().String(),
														 uPIN,pToken,userID,sID,priKey)
		if err != nil {
				log.Fatal(err)
		}
		fmt.Println(string(bt))
	}
	if cmd == "9" {
		priKey, _, sID, userID, _ := GetWalletInfo()
		data, err := mixin.NetworkSnapshots(mixin.GetAssetId("BTC"), time.Now().AddDate(0, 0, -1), true, 10, userID, sID, priKey)
		if err != nil { log.Fatal(err) }
		log.Println(string(data))
	}
  if cmd == "a" {
    QueryInfo, err := mixin.VerifyPIN(PinCode, PinToken,UserId,
                                      SessionId,PrivateKey)
    if err != nil {
            log.Fatal(err)
    }
    fmt.Println(string(QueryInfo))
  }
  if cmd == "v" {
    priKey, pinTkn, sID, userID, pinX := GetWalletInfo()
    QueryInfo, err := mixin.VerifyPIN(pinX, pinTkn,userID,
                                      sID,priKey)
    if err != nil {
            log.Fatal(err)
    }
    fmt.Println(string(QueryInfo))
  }
	}
  // c6d0c728-2624-429b-8e0d-d9d19b6592fa
}
