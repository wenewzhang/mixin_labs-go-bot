# How to trade bitcoin through Golang

## Solution One: pay to ExinCore API
[Exincore](https://github.com/exinone/exincore) provide a commercial trading API on Mixin Network.

You pay USDT to ExinCore, ExinCore transfer Bitcoin to you on the fly with very low fee and fair price. Every transaction is anonymous to public but still can be verified on blockchain explorer. Only you and ExinCore know the details.

ExinCore don't know who you are because ExinCore only know your client's uuid.

### Pre-request:
You should  have created a bot based on Mixin Network. Create one by reading [GO Bitcoin tutorial](https://github.com/wenewzhang/mixin_labs-go-bot).

#### Install required packages
As you know, we introduce you the **mixin-sdk-go** in [chapter 1](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README.md), assume it has installed before, let's install **uuid, msgpack** here.
```bash
  go get -u github.com/vmihailenco/msgpack
  go get -u github.com/satori/go.uuid
```
#### Deposit USDT or Bitcoin into your Mixin Network account and read balance
ExinCore can exchange between Bitcoin, USDT, EOS, Eth etc. Here show you how to exchange between USDT and Bitcoin,
Check the wallet's balance & address before you make order.

- Check the address & balance, remember it Bitcoin wallet address.
- Deposit Bitcoin to this Bitcoin wallet address.
- Check Bitcoin balance after 100 minutes later.

**By the way, Bitcoin & USDT 's address are the same.**

```go
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
```
#### Read market price
How to check the coin's price? You need understand what is the base coin. If you want buy Bitcoin and sell USDT, the USDT is the base coin. If you want buy USDT and sell Bitcoin, the Bitcoin is the base coin.
```go
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
```

#### Create a memo to prepare order
The chapter two: [Echo Bitcoin](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README2.md) introduce transfer coins. But you need to let ExinCore know which coin you want to buy. Just write your target asset into memo.
```go
packUuid, _ := uuid.FromString(mixin.GetAssetId("BTC"))
pack, _ := msgpack.Marshal(OrderAction{A: packUuid})
memo := base64.StdEncoding.EncodeToString(pack)
// fmt.Println(memo)
```
#### Pay BTC to API gateway with generated memo
Transfer Bitcoin(BTC_ASSET_ID) to ExinCore(EXIN_BOT), put you target asset uuid in the memo, otherwise, ExinCore will refund you coin immediately!
```go
EXIN_BOT        = "61103d28-3ac2-44a2-ae34-bd956070dab1";

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
```
The ExinCore should transfer the target coin to your bot, meanwhile, put the fee, order id, price etc. information in the memo, unpack the data like below.
- **readUserSnapshots** Read snapshots of the user.
```go
type OrderResponse struct {
    C  int    // code
    P  string     // price, only type is return
    F  string     // ExinCore fee, only type is return
    FA []byte     // ExinCore fee asset, only type is return
    T  string     // type: refund(F)|return(R)|Error(E)
    O  uuid.UUID  // order: trace_id
}
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
```

If you coin exchange successful, console output like below:
```bash
---------------Successful----Exchange-------------
You got  0.3981012
815b0b1a-2764-3736-8faa-42d694fa620a  price: 3996.8  Fee: 0.0007994
```

#### Read Bitcoin balance
Check the wallet's balance.
```go
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
```
## Source code usage
Execute **go run coin_exchange.go** to run it.

- 1: Create Wallet
- 2: Read Bitcoin balance & Address
- 3: Read USDT balance & Address
- 4: Read EOS balance & address
- 5: pay 0.0001 BTC buy USDT
- 6: Read ExinCore Price(USDT)
- 7: Read ExinCore Price(BTC)
- 8: pay 1 USDT buy BTC
- q: Exit
Make your choose:

[Full source code](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/coin_exchange/coin_exchange.go)

## Solution Two: List your order on Ocean.One exchange
