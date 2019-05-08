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
    "strconv"
    mixin "github.com/MooooonStar/mixin-sdk-go/network"
		"github.com/MooooonStar/mixin-sdk-go/messenger"
)

const (
	UserId     = "21042518-85c7-4903-bb19-f311813d1f51"
	PinCode    = "152997"
	SessionId  = "f8c55131-f78a-4858-9ec1-7c69d2a88d0d"
	PinToken   = "aUO0NSHchcqGony2gSDJqW2ohStqF47nJlAeAo6dFgXSD0cg/RaJtKT6fRjN63q7wZGNYwOOTzPIq6WypnrZ1AR0spYE8dZ6thWAooIM2alVwGtjofczVdvPeOegCCbDgcIGTTxTKgAPij10AHaI2RX1Is4gB3zzArF5c8l54no="
	//please delele the blank of PrivateKey the before each line
	PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDFOoiKwycPZCtM/kjBbuNbT3lP2eKfV4lTHKRj2UEfvs7RMRhk
7hzY0bxhLPP6ZI64RI6EAobRUUiK8MB4pqKzO0AukUUCwizrIN0LsWvH+qZXIujf
JAGsdh95UQBk5cJWG68xg1EdijDNDEoa+DXoDnWdNjQJSRBVx2D32UHHTQIDAQAB
AoGAfp5Xbo5fEziBvAo790MTX1mkTilZnmZ6WQs4Vonxj0nWSOK2AIYFqwTrZY+Q
ip3oKlCJFiLxHoyKf/iT+GEybbEBEwTwum+I/NQA+dQixLxoBP5jnSrt9HoPxJ2h
sheoBfI/OT1+0QqDlqVlnJeNREkmLHXqdU2r2V5RXBeg0FkCQQD0aYqXCBfb4mln
Kt9xx8o+yZobzGqgti5IW8Nw+OONGMlkTt+eXiQp43oPH9PmznbNEhnWAXIIG+Hi
PrHb1BbjAkEAzpRRs9qSPABF/mQx1u9AVwmGKtLVzQ6HldhChCzwWzFYqu6wReMm
r4Gn+LKSJNOjZsgII7AFlu8tXMBGDpTQDwJBAIvTWXMgMS4dcHmSIHTifMTA50Zi
Atpgf0fsH3qhGOVeudCGAw6CAyRnvCus5LiVg4e8hEVXXFphQTAC+BOwWUsCQQC/
7NzblD44sKhW6Q/E+RN1yct1DdzFXrJpbTqfQoEsuHQAmzH6PEg81uEQFhfhTx+I
5l9piKgoyp4ChkCQW4HRAkB8bQ/UkC9iny345GoCoy/Pf6oSfSttokFP7Z9vaERH
FaFESfvfy05ogzB5hN3LoywwSLymrHgeQQK2RYunSpAS
-----END RSA PRIVATE KEY-----`
	EXIN_BOT                = "61103d28-3ac2-44a2-ae34-bd956070dab1"
  OCEANONE_BOT            = "aaff5bef-42fb-4c9f-90e0-29f69176b7d4";
  MASTER_UUID             = "0b4f49dc-8fb4-4539-9a89-fb3afc613747";
  ERC20_BENZ              = "2b9c216c-ef60-398d-a42a-eba1b298581d";
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
type OceanOrderAction struct {
  S string    // side
  A uuid.UUID // asset
  P string    // price
  T string    // type
}
type OceanOrderCancel struct {
  O uuid.UUID // asset
}
type OrderResponse struct {
    C  int    // code
    P  string     // price, only type is return
    F  string     // ExinCore fee, only type is return
    FA []byte     // ExinCore fee asset, only type is return
    T  string     // type: refund(F)|return(R)|Error(E)
    O  uuid.UUID  // order: trace_id
}

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

func ReadAssetBalance(asset_id, UserID, SessionID, PrivateKey string) (string) {
  var AssetInfo map[string]interface{}
  AssetInfoBytes, err  := mixin.ReadAsset(mixin.GetAssetId(asset_id),
                                         UserID,SessionID,PrivateKey)
  if err != nil { log.Fatal(err) }
  // fmt.Println(string(AssetInfoBytes))
  if err := json.Unmarshal(AssetInfoBytes, &AssetInfo); err != nil {
      log.Fatal(err)
  }
  return AssetInfo["data"].(map[string]interface{})["balance"].(string)
}

func ReadAssetBalanceByUUID(asset_uuid, UserID, SessionID, PrivateKey string) (string) {
  var AssetInfo map[string]interface{}
  AssetInfoBytes, err  := mixin.ReadAsset(asset_uuid,
                                         UserID,SessionID,PrivateKey)
  if err != nil { log.Fatal(err) }
  // fmt.Println(string(AssetInfoBytes))
  if err := json.Unmarshal(AssetInfoBytes, &AssetInfo); err != nil {
      log.Fatal(err)
  }
  return AssetInfo["data"].(map[string]interface{})["balance"].(string)
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

func FormatOceanOneMarketPrice(asset_id string, base_asset string) {
  priceInfo, err := GetOceanOneMarketPrice(asset_id, base_asset)
  if err != nil {
    log.Fatal(err)
  }

  var marketInfo map[string]interface{}
  err = json.Unmarshal([]byte(priceInfo), &marketInfo)
   fmt.Println("Price | Amount | Funds | Side")
  for _, v := range (marketInfo["data"].
                    (map[string]interface{})["data"].
                    (map[string]interface{})["asks"].
                    ([]interface{})) {
    fmt.Println(v.(map[string]interface{})["price"],
                v.(map[string]interface{})["amount"],
                v.(map[string]interface{})["funds"],
                v.(map[string]interface{})["side"],
               )
  }
  for _, v := range (marketInfo["data"].
                    (map[string]interface{})["data"].
                    (map[string]interface{})["bids"].
                    ([]interface{})) {
    fmt.Println(v.(map[string]interface{})["price"],
                v.(map[string]interface{})["amount"],
                v.(map[string]interface{})["funds"],
                v.(map[string]interface{})["side"],
               )
  }
}
func GetOceanOneMarketPrice(asset_id string, base_asset string) ([]byte, error)  {
	var body []byte
	req, err := http.NewRequest("GET", "https://events.ocean.one/markets/" + asset_id + "-" + base_asset + "/book",
                              bytes.NewReader(body))
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

func transferBotWallet(AssetID,opponentID,PinCode,PinToken,UserId,SessionId,PrivateKey string) {
  balance := ReadAssetBalance(AssetID,UserId,SessionId,PrivateKey)
  fmt.Println(balance)
  ibalance, _ := strconv.ParseFloat(balance,64)
  if ibalance > 0 {
    transInfo, _ := mixin.Transfer(opponentID,balance,mixin.GetAssetId(AssetID),"memo",
                             messenger.UuidNewV4().String(),
                             PinCode,PinToken,UserId,SessionId,PrivateKey)
    fmt.Println(string(transInfo))
  }
}

func transferBotWalletByUUID(AssetUUID,opponentID,PinCode,PinToken,UserId,SessionId,PrivateKey string) {
  balance := ReadAssetBalanceByUUID(AssetUUID,UserId,SessionId,PrivateKey)
  fmt.Println(balance)
  ibalance, _ := strconv.ParseFloat(balance,64)
  if ibalance > 0 {
    transInfo, _ := mixin.Transfer(opponentID,balance,AssetUUID,"memo",
                             messenger.UuidNewV4().String(),
                             PinCode,PinToken,UserId,SessionId,PrivateKey)
    fmt.Println(string(transInfo))
  }
}

func generateOceanOrderMemo(TargetAsset, Side, Price string) (string) {
  packUuid, _ := uuid.FromString(TargetAsset)
  memoOcean,_ :=
    msgpack.Marshal(OceanOrderAction{
      T: "L",
      P: Price,
      S: Side,
      A: packUuid,
    })
  return  base64.StdEncoding.EncodeToString(memoOcean)
}

func main() {
  // Pack memo
  packUuid, _ := uuid.FromString("c6d0c728-2624-429b-8e0d-d9d19b6592fa")
  fmt.Println(packUuid)
  pack, _ := msgpack.Marshal(OrderAction{A: packUuid,})
  fmt.Println(pack)
  memo := base64.StdEncoding.EncodeToString(pack)
  // gaFBxBDG0McoJiRCm44N2dGbZZL6
  fmt.Println(memo)
  // Parse memo
  parsedpack, _ := base64.StdEncoding.DecodeString(memo)
  orderAction := OrderAction{}
  _ = msgpack.Unmarshal(parsedpack, &orderAction)
  fmt.Println(orderAction.A)
  memoOcean,_ :=
    msgpack.Marshal(OceanOrderAction{
      T: "L",
      P: "0.1",
      S: "A",
      A: packUuid,
    })
  memoOceanB64 := base64.StdEncoding.EncodeToString(memoOcean)
  fmt.Println(memoOceanB64)

  scanner   := bufio.NewScanner(os.Stdin)
	var PromptMsg string
	PromptMsg  = "1: Create Wallet\n2: Read Bitcoin balance & Address \n3: Read USDT balance & Address\n4: Read EOS balance & address\n"
  PromptMsg += "tbb:Transfer BTC from Bot to Wallet\ntbm:Transfer BTC from Wallet to Master\n"
  PromptMsg += "teb:Transfer EOS from Bot to Wallet\ntem:Transfer EOS from Wallet to Master\n"
  PromptMsg += "tub:Transfer USDT from Bot to Wallet\ntum:Transfer USDT from Wallet to Master\n"
  PromptMsg += "tcb:Transfer CNB from Bot to Wallet\ntcm:Transfer CNB from Wallet to Master\n"
  PromptMsg += "txb:Transfer XIN from Bot to Wallet\ntxm:Transfer XIN from Wallet to Master\n"
  PromptMsg += "trb:Transfer ERC20 from Bot to Wallet\ntrm:Transfer ERC20 from Wallet to Master\n"
  PromptMsg += "5: pay 0.0001 BTC buy USDT\n6: Read ExinCore Price(USDT)\n7: Read ExinCore Price(BTC)\n"
	PromptMsg += "8: pay 1 USDT buy BTC\n9: Read Snapshots\na: Verify bot PIN code\nv: Verify wallet PIN code\n"
  PromptMsg += "ab: Read Bot Assets\naw: Read Wallet Assets\n";
	PromptMsg += "o: Ocean.One Exchange\nq: Exit \nMake your choose:"
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
    fmt.Println(userInfo["data"])
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
		for _, v := range (marketInfo["data"].([]interface{})) {
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
		for _, v := range (marketInfo["data"].([]interface{})) {
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
    fmt.Println("Input the trade time:")
    var tmUTC string
    scanner.Scan()
    tmUTC = scanner.Text()
    tm, _:= time.Parse(time.RFC3339Nano,tmUTC)
		snapData, err := mixin.NetworkSnapshots("", tm, true, 3, userID, sID, priKey)
		if err != nil { log.Fatal(err) }
		fmt.Println(string(snapData))
    // fmt.Println(snapData.data)
    var snapInfo map[string]interface{}
    err = json.Unmarshal([]byte(snapData), &snapInfo)
    if err != nil {
        log.Fatal(err)
    }
    for _, v := range (snapInfo["data"].([]interface{})) {
      val := v.(map[string]interface{})["amount"]
      if amount, ok := val.(string); ok {
          if v.(map[string]interface{})["data"] != nil {
            strMemo := v.(map[string]interface{})["data"]
            memo := strMemo.(string)
            parsedpack, _ := base64.StdEncoding.DecodeString(memo)
            orderAction := OrderResponse{}
            _ = msgpack.Unmarshal(parsedpack, &orderAction)
            if orderAction.C == 1000 {
              fmt.Println("---------------Successful----Exchange-------------")
              fmt.Println("You got ",amount)
              uuidAsset,_:= uuid.FromBytes(orderAction.FA)
              fmt.Println(uuidAsset," price:",orderAction.P," Fee:",orderAction.F)
            }
          }
      } else {
          continue
      }
    }
	}
  if cmd == "a" {
    QueryInfo, err := mixin.VerifyPIN(PinCode, PinToken,UserId,
                                      SessionId,PrivateKey)
    if err != nil {
            log.Fatal(err)
    }
    fmt.Println(string(QueryInfo))
  }
  if cmd == "aw" {
    priKey, _, sID, userID, _ := GetWalletInfo()
    assets, err := mixin.ReadAssets(userID,sID,priKey)
    if err != nil {
        log.Fatal(err)
    }
    var AssetsInfo map[string]interface{}
    err = json.Unmarshal(assets, &AssetsInfo)
    if err != nil {
        log.Fatal(err)
    }
    // fmt.Println("Data is: ",AssetsInfo["data"].(map[string]interface{})["public_key"])
    for _, v := range (AssetsInfo["data"].([]interface{})) {
      if v.(map[string]interface{})["symbol"] == "EOS" {
        fmt.Println(v.(map[string]interface{})["symbol"]," ",
                    v.(map[string]interface{})["account_name"]," ",
                    v.(map[string]interface{})["account_tag"]," ",
                    v.(map[string]interface{})["balance"])
      } else {
        fmt.Println(v.(map[string]interface{})["symbol"]," ",
                    v.(map[string]interface{})["public_key"]," ",
                    v.(map[string]interface{})["balance"])
      }
    }
  }
  if cmd == "ab" {
    // priKey, _, sID, userID, _ := GetWalletInfo()
    assets, err := mixin.ReadAssets(UserId,SessionId,PrivateKey)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(string(assets))
    var AssetsInfo map[string]interface{}
    err = json.Unmarshal(assets, &AssetsInfo)
    if err != nil {
        log.Fatal(err)
    }
    // fmt.Println("Data is: ",AssetsInfo["data"].(map[string]interface{})["public_key"])
    for _, v := range (AssetsInfo["data"].([]interface{})) {
      if v.(map[string]interface{})["symbol"] == "EOS" {
        fmt.Println(v.(map[string]interface{})["symbol"]," ",
                    v.(map[string]interface{})["account_name"]," ",
                    v.(map[string]interface{})["account_tag"]," ",
                    v.(map[string]interface{})["balance"])
      } else {
        fmt.Println(v.(map[string]interface{})["symbol"]," ",
                    v.(map[string]interface{})["public_key"]," ",
                    v.(map[string]interface{})["balance"])
      }
    }
  }

  if cmd == "tcb" {
    _, _, _, userID, _ := GetWalletInfo()
    transferBotWallet("CNB",userID,PinCode,PinToken,UserId,SessionId,PrivateKey)
  }
  if cmd == "tcm" {
    priKey, pToken, sID, userID, uPIN := GetWalletInfo()
    transferBotWallet("CNB",MASTER_UUID,uPIN,pToken,userID,sID,priKey)
  }
  if cmd == "tub" {
    _, _, _, userID, _ := GetWalletInfo()
    transferBotWallet("USDT",userID,PinCode,PinToken,UserId,SessionId,PrivateKey)
  }
  if cmd == "tum" {
    priKey, pToken, sID, userID, uPIN := GetWalletInfo()
    transferBotWallet("USDT",MASTER_UUID,uPIN,pToken,userID,sID,priKey)
  }
  if cmd == "txb" {
    _, _, _, userID, _ := GetWalletInfo()
    transferBotWallet("XIN",userID,PinCode,PinToken,UserId,SessionId,PrivateKey)
  }
  if cmd == "txm" {
    priKey, pToken, sID, userID, uPIN := GetWalletInfo()
    transferBotWallet("XIN",MASTER_UUID,uPIN,pToken,userID,sID,priKey)
  }
  if cmd == "tbb" {
    _, _, _, userID, _ := GetWalletInfo()
    transferBotWallet("BTC",userID,PinCode,PinToken,UserId,SessionId,PrivateKey)
  }
  if cmd == "tbm" {
    priKey, pToken, sID, userID, uPIN := GetWalletInfo()
    transferBotWallet("BTC",MASTER_UUID,uPIN,pToken,userID,sID,priKey)
  }
  if cmd == "teb" {
    _, _, _, userID, _ := GetWalletInfo()
    transferBotWallet("EOS",userID,PinCode,PinToken,UserId,SessionId,PrivateKey)
  }
  if cmd == "tem" {
    priKey, pToken, sID, userID, uPIN := GetWalletInfo()
    transferBotWallet("EOS",MASTER_UUID,uPIN,pToken,userID,sID,priKey)
  }
  if cmd == "trb" {
    _, _, _, userID, _ := GetWalletInfo()
    transferBotWalletByUUID(ERC20_BENZ,userID,PinCode,PinToken,UserId,SessionId,PrivateKey)
  }
  if cmd == "trm" {
    priKey, pToken, sID, userID, uPIN := GetWalletInfo()
    transferBotWalletByUUID(ERC20_BENZ,MASTER_UUID,uPIN,pToken,userID,sID,priKey)
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
  if cmd == "o" {
    scanner   := bufio.NewScanner(os.Stdin)
    var PromptMsg string
    PromptMsg  = "1:  Fetch XIN/USDT orders\ns1: Sell XIN/USDT\nb1: Buy XIN/USDT\n"
    PromptMsg += "2:  Fetch ERC20(Benz)/USDT orders\ns2: Sell Benz/USDT\nb2: Buy Benz/USDT\n"
    PromptMsg += "c: Cancel Order\nq:  Exit\n"
    for  {
     fmt.Print(PromptMsg)
     var cmd string
     scanner.Scan()
     cmd = scanner.Text()
     if cmd == "q" { break }
     if cmd == "1" {
     		FormatOceanOneMarketPrice(mixin.GetAssetId("XIN"),mixin.GetAssetId("USDT"))
     	}
     if cmd == "2" {
        FormatOceanOneMarketPrice(ERC20_BENZ,mixin.GetAssetId("USDT"))
      }
     if cmd == "s1" {
       fmt.Print("Please input the price of XIN/USDT: ")
       var pcmd string
       var acmd string
       scanner.Scan()
       pcmd = scanner.Text()
       fmt.Println(pcmd)
       fmt.Print("Please input the amount of XIN: ")
       scanner.Scan()
       acmd = scanner.Text()
       fmt.Println(acmd)
       omemo := generateOceanOrderMemo(mixin.GetAssetId("USDT"),"A",pcmd)
       priKey, pToken, sID, userID, uPIN := GetWalletInfo()
       balance := ReadAssetBalance("XIN",userID,sID,priKey)
       fmt.Println(balance)
       fbalance, _ := strconv.ParseFloat(balance,64)
       abalance, _ := strconv.ParseFloat(acmd,64)
       if fbalance > 0 && fbalance >= abalance {
         fmt.Println(omemo)
         transInfo, _ := mixin.Transfer(OCEANONE_BOT,
                                        acmd,
                                        mixin.GetAssetId("XIN"),
                                        omemo,
                                        messenger.UuidNewV4().String(),
                                        uPIN,pToken,userID,sID,priKey)
         fmt.Println(string(transInfo))
       }
     }
     if cmd == "b1" {
       fmt.Print("Please input the price of XIN/USDT: ")
       var pcmd string
       var acmd string
       scanner.Scan()
       pcmd = scanner.Text()
       fmt.Println(pcmd)
       fmt.Print("Please input the amount of USDT: ")
       scanner.Scan()
       acmd = scanner.Text()
       fmt.Println(acmd)
       omemo := generateOceanOrderMemo(mixin.GetAssetId("XIN"),"B",pcmd)
       priKey, pToken, sID, userID, uPIN := GetWalletInfo()
       balance := ReadAssetBalance("USDT",userID,sID,priKey)
       fmt.Println(balance)
       fbalance, _ := strconv.ParseFloat(balance,64)
       abalance, _ := strconv.ParseFloat(acmd,64)
       if fbalance > 0 && fbalance >= abalance {
         fmt.Println(omemo)
         transInfo, _ := mixin.Transfer(OCEANONE_BOT,
                                        acmd,
                                        mixin.GetAssetId("USDT"),
                                        omemo,
                                        messenger.UuidNewV4().String(),
                                        uPIN,pToken,userID,sID,priKey)
         fmt.Println(string(transInfo))
       }
     }//end of b1
     if cmd == "c" {
       fmt.Print("Please input the Order id: ")
       var ocmd string
       scanner.Scan()
       ocmd = scanner.Text()
       fmt.Println(ocmd)
       memoOcean,_ :=
         msgpack.Marshal(OceanOrderCancel{
           O: packUuid,
         })
       omemoCancel := base64.StdEncoding.EncodeToString(memoOcean)
       priKey, pToken, sID, userID, uPIN := GetWalletInfo()
       balance := ReadAssetBalance("CNB",userID,sID,priKey)
       fmt.Println(balance)
       fbalance, _ := strconv.ParseFloat(balance,64)
       // abalance, _ := strconv.ParseFloat(acmd,64)
       if fbalance > 0 && fbalance >= 0.0000001 {
         fmt.Println(omemoCancel)
         transInfo, _ := mixin.Transfer(OCEANONE_BOT,
                                        "0.00000001",
                                        mixin.GetAssetId("CNB"),
                                        omemoCancel,
                                        messenger.UuidNewV4().String(),
                                        uPIN,pToken,userID,sID,priKey)
         fmt.Println(string(transInfo))
       } else { fmt.Println("Not enough CNB!") }
     }
    }
   }//end of Ocean.one exchange
  }
  // c6d0c728-2624-429b-8e0d-d9d19b6592fa
}
