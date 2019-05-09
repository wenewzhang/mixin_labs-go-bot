# How to list any ERC20 token on decentralized market through Go
![cover](https://github.com/wenewzhang/mixin_labs-go-bot/raw/master/Bitcoin_go.jpg)

OceanOne is introduced in [last chapter](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README5.md), you can trade Bitcoin. All kinds of crypto asset on Mixin Network can be listed on OceanOne.All ERC20 token and EOS token can be listed. Following example will show you how to list a ERC20 token.

There is a [ERC20 token](https://etherscan.io/token/0xc409b5696c5f9612e194a582e14c8cd41ecdbc67) called Benz. It is deposited into Mixin Network. You can search all transaction history from [Mixin Network browser](https://mixin.one/snapshots/2b9c216c-ef60-398d-a42a-eba1b298581d )

### Pre-request:
Deposit some coin to your wallet, and then use **ReadAssets** API fetch the asset UUID which Mixin Network gave it.

### Get the ERC-20 compliant coin UUID
The **ReadAssets** API return json data, for example:

- **asset_id** UUID of this coin
- **public_key** The wallet address for this coin
- **symbol**  Coin name, Eg: Benz.

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
                  v.(map[string]interface{})["asset_id"]," ",
                  v.(map[string]interface{})["account_name"]," ",
                  v.(map[string]interface{})["account_tag"]," ",
                  v.(map[string]interface{})["balance"])
    } else {
      fmt.Println(v.(map[string]interface{})["symbol"]," ",
                  v.(map[string]interface{})["asset_id"]," ",
                  v.(map[string]interface{})["public_key"]," ",
                  v.(map[string]interface{})["balance"])
    }
  }
}
```

The detail information of **ReadAssets** is output like below:

```bash
Make your choose:aw
Benz   0x330860Ec473fF366F5Bc4339a69f5bffB52d18Cb   88.9
EOS   eoswithmixin   79dd76cedf8f6af49a8d98216bbde890   0
USDT   16wWhKAjmACvZzkfxkyrVutqfrJ1JQ83aj   1
CNB   0x330860Ec473fF366F5Bc4339a69f5bffB52d18Cb   0.10999989
BTC   16wWhKAjmACvZzkfxkyrVutqfrJ1JQ83aj   0
XIN   0x330860Ec473fF366F5Bc4339a69f5bffB52d18Cb   0.01
```
### Make the limit order
- **Limit Order to Buy**  at or below the market.
- **Limit Order to Sell**  at or above the market.

OceanOne support three base coin: USDT, XIN, BTC, that mean you can sell or buy it between USDT, XIN, BTC, so, you have there order: Benz/USDT, Benz/XIN, Benz/BTC, here show you how to make the sell order with USDT.

### Make the limit order to sell.

```go
if cmd == "s2" {
  fmt.Print("Please input the price of ERC20/USDT: ")
  var pcmd string
  var acmd string
  scanner.Scan()
  pcmd = scanner.Text()
  fmt.Println(pcmd)
  fmt.Print("Please input the amount of ERC20: ")
  scanner.Scan()
  acmd = scanner.Text()
  fmt.Println(acmd)
  omemo := generateOceanOrderMemo(mixin.GetAssetId("USDT"),"A",pcmd)
  priKey, pToken, sID, userID, uPIN := GetWalletInfo()
  balance := ReadAssetBalanceByUUID(ERC20_BENZ,userID,sID,priKey)
  fmt.Println(balance)
  fbalance, _ := strconv.ParseFloat(balance,64)
  abalance, _ := strconv.ParseFloat(acmd,64)
  if fbalance > 0 && fbalance >= abalance {
    fmt.Println(omemo)
    transInfo, _ := mixin.Transfer(OCEANONE_BOT,
                                   acmd,
                                   ERC20_BENZ,
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
  } else { fmt.Println("Not enough BenZ!") }
}
```

### Make the limit order to buy.
After the order commit, wait 1 minute to let the OceanOne exchange initialize it.
```go
if cmd == "b2" {
  fmt.Print("Please input the price of ERC20/USDT: ")
  var pcmd string
  var acmd string
  scanner.Scan()
  pcmd = scanner.Text()
  fmt.Println(pcmd)
  fmt.Print("Please input the amount of USDT: ")
  scanner.Scan()
  acmd = scanner.Text()
  fmt.Println(acmd)
  omemo := generateOceanOrderMemo(ERC20_BENZ,"B",pcmd)
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
}//end of b2
```

### Read orders book from Ocean.one
Now, check the orders-book.

```go
if cmd == "3" {
   FormatOceanOneMarketPrice(mixin.GetAssetId("BTC"),mixin.GetAssetId("USDT"))
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

### Command of make orders

Commands list of this source code:

- trb:Transfer ERC20 from Bot to Wallet
- trm:Transfer ERC20 from Wallet to Master
- o: Ocean.One Exchange

Make your choose(eg: q for Exit!):
- 3:  Orders-Book of ERC20/USDT
- b3: Buy ERC20 pay USDT
- s3: Sell ERC20 get USDT
- c: Cancel the order
- q: Exit

[Full source code](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/coin_exchange/coin_exchange.go)
