# How to list bitcoin order through Java
![cover](https://github.com/wenewzhang/mixin_labs-go-bot/raw/master/Bitcoin_go.jpg)

Exincore is introduced in [last chapter](https://github.com/wenewzhang/mixin_labs-java-bot/blob/master/README4.md), you can exchange many crypto asset at market price and receive your asset in 1 seconds. If you want to trade asset at limited price, or trade asset is not supported by ExinCore now, OceanOne is the answer.
## Solution Two: List your order on Ocean.One exchange
[Ocean.one](https://github.com/mixinNetwork/ocean.one) is a decentralized exchange built on Mixin Network, it's almost the first time that a decentralized exchange gain the same user experience as a centralized one.

You can list any asset on OceanOne. Pay the asset you want to sell to OceanOne account, write your request in payment memo, OceanOne will list your order to market. It send asset to your wallet after your order is matched.

* No sign up required
* No deposit required
* No listing process.

### Pre-request:
You should  have created a bot based on Mixin Network. Create one by reading [Java Bitcoin tutorial](https://github.com/wenewzhang/mixin_labs-java-bot/blob/master/README.md).

#### Install required packages
[Chapter 4](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README4.md), assume it has installed before.

#### Deposit USDT or Bitcoin into your Mixin Network account and read balance
The Ocean.one can match any order. Here we exchange between USDT and Bitcoin, Check the wallet's balance & address before you make order.

- Check the address & balance, find it's Bitcoin wallet address.
- Deposit Bitcoin to this Bitcoin wallet address.
- Check Bitcoin balance after 100 minutes later.

**Omni USDT address is same as Bitcoin address**

```go
userInfo, userID := ReadAssetInfo("USDT")
fmt.Println("User ID ",userID, "'s USDT Address is: ",
           userInfo["data"].(map[string]interface{})["public_key"])
fmt.Println("Balance is: ",
           userInfo["data"].(map[string]interface{})["balance"])

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
```

#### Read orders book from Ocean.one
How to check the coin's price? You need understand what is the base coin. If you want buy Bitcoin and sell USDT, the USDT is the base coin. If you want buy USDT and sell Bitcoin, the Bitcoin is the base coin.


```go
if cmd == "1" {
   FormatOceanOneMarketPrice(mixin.GetAssetId("XIN"),mixin.GetAssetId("USDT"))
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
```

#### Create a memo to prepare order
The chapter two: [Echo Bitcoin](https://github.com/wenewzhang/mixin_labs-java-bot/blob/master/README2.md) introduce transfer coins. But you need to let Ocean.one know which coin you want to buy.
- **Side** "B" or "A", "B" for buy, "A" for Sell.
- **AssetUUID** UUID of the asset you want to buy
- **Price** If Side is "B", Price is AssetUUID; if Side is "A", Price is the asset which transfer to Ocean.one.

```go
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
```

#### Pay XIN to OceanOne with generated memo
Transfer XIN to Ocean.one(OCEANONE_BOT), put you target asset uuid(USDT) in the memo.
```go
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
    var jsTransInfo map[string]interface{}
    err := json.Unmarshal([]byte(transInfo), &jsTransInfo)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("The Order id is " + jsTransInfo["data"].(map[string]interface{})["trace_id"].(string) +
               " it is needed to cancel the order!")
  } else { fmt.Println("Not enough XIN!") }
}
```

If you want sell USDT buy XIN, call it like below:
```go
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
    var jsTransInfo map[string]interface{}
    err := json.Unmarshal([]byte(transInfo), &jsTransInfo)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("The Order id is " + jsTransInfo["data"].(map[string]interface{})["trace_id"].(string) +
               " it is needed to cancel the order!")
  } else { fmt.Println("Not enough USDT!") }
}//end of b1
```

A success order output like below:

```bash
Please input the price of BTC/USDT: 5666
5666
Please input the amount of USDT: 1
1
1
hKFToUKhQcQQxtDHKCYkQpuODdnRm2WS+qFQpDU2NjahVKFM
{"data":{"type":"transfer","snapshot_id":"c1518f2c-c2e8-4d2d-b7c3-a42770f2bdab",
"opponent_id":"aaff5bef-42fb-4c9f-90e0-29f69176b7d4",
"asset_id":"815b0b1a-2764-3736-8faa-42d694fa620a","amount":"-1",
"trace_id":"26ef0ec1-a60c-4702-b8ca-4bd40330e120",
"memo":"hKFToUKhQcQQxtDHKCYkQpuODdnRm2WS+qFQpDU2NjahVKFM",
"created_at":"2019-05-08T06:48:42.919216755Z",
"counter_user_id":"aaff5bef-42fb-4c9f-90e0-29f69176b7d4"}}
The Order id is 26ef0ec1-a60c-4702-b8ca-4bd40330e120 it is needed to cancel the order!
```

## Cancel the Order
To cancel order, just pay any amount of any asset to OceanOne, and write trace_id into memo. Ocean.one take the trace_id as the order id, for example, **26ef0ec1-a60c-4702-b8ca-4bd40330e120** is a order id,
We can use it to cancel the order.

```go
if cmd == "c" {
  fmt.Print("Please input the Order id: ")
  var ocmd string
  scanner.Scan()
  ocmd = scanner.Text()
  fmt.Println(ocmd)
  orderid, _ := uuid.FromString(ocmd)
  memoOcean,_ :=
    msgpack.Marshal(OceanOrderCancel{
      O: orderid,
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
```

#### Read Bitcoin balance
Check the wallet's balance.
```go
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
```

## Source code usage
Build it and then run it.

- **go run coin_exchange.go**  build project.

Commands list of this source code:

- aw: Read Wallet Assets
- o: Ocean.One Exchange
- q: Exit
Make your choose:

Make your choose(eg: q for Exit!):

- 1:  Fetch XIN/USDT orders
- s1: Sell XIN/USDT
- b1: Buy XIN/USDT
- 2:  Fetch ERC20(Benz)/USDT orders
- s2: Sell Benz/USDT
- b2: Buy Benz/USDT
- 3:  Fetch BTC/USDT orders
- s3: Sell BTC/USDT
- b3: Buy BTC/USDT
- c: Cancel Order
- q:  Exit

[Full source code](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/coin_exchange/coin_exchange.go)
