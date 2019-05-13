# 通过 Go 在去中心化交易所OceanOne上挂单买卖任意ERC20 token
![cover](https://github.com/wenewzhang/mixin_labs-go-bot/raw/master/Bitcoin_go.jpg)

在[上一课](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README5.md)中，我们介绍了如何在OceanOne交易比特币。OceanOne支持交易任何Mixin Network上的token，包括所有的ERC20和EOS token，不需要任何手续和费用，直接挂单即可。下面介绍如何将将一个ERC20 token挂上OceanOne交易！在掌握了ERC20 token之后，就可以把任何token在Ocean上买卖。

此处我们用一个叫做Benz的[ERC20 token](https://etherscan.io/token/0xc409b5696c5f9612e194a582e14c8cd41ecdbc67)为例。这个token已经被充值进Mixin Network，你可以在[区块链浏览器](https://mixin.one/snapshots/2b9c216c-ef60-398d-a42a-eba1b298581d )看到这个token在Mixin Network内部的总数和交易
### 预备知识:
先将Ben币存入你的钱包，然后使用**ReadAssets** API读取它的UUID.

### 取得该币的UUID
调用 **ReadAssets** API 会返回json数据, 如:

- **asset_id** 币的UUID.
- **public_key** 该币的当前钱包的地址.
- **symbol**  币的名称. 如: Benz.

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

调用 **ReadAssets** API的完整输出如下:

```bash
Make your choose:aw
Benz   0x330860Ec473fF366F5Bc4339a69f5bffB52d18Cb   88.9
EOS   eoswithmixin   79dd76cedf8f6af49a8d98216bbde890   0
USDT   16wWhKAjmACvZzkfxkyrVutqfrJ1JQ83aj   1
CNB   0x330860Ec473fF366F5Bc4339a69f5bffB52d18Cb   0.10999989
BTC   16wWhKAjmACvZzkfxkyrVutqfrJ1JQ83aj   0
XIN   0x330860Ec473fF366F5Bc4339a69f5bffB52d18Cb   0.01
```

### 限价挂单
- **挂限价买单**  低于或者等于市场价的单.
- **挂限价卖单**  高于或者是等于市场价的单.

OceanOne支持三种基类价格: USDT, XIN, BTC, 即: Benz/USDT, Benz/XIN, Benz/BTC, 这儿示范Benz/USDT.

### 限价挂卖单.
新币挂单后,需要等一分钟左右，等OceanOne来初始化新币的相关数据.

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

### 限价挂买单.
新币挂单后,需要等一分钟左右，等OceanOne来初始化新币的相关数据.

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
### 读取币的价格列表
读取币的价格列表，来确认挂单是否成功!

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
### ERC20相关的操作指令

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

[完整的代码](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/coin_exchange/coin_exchange.go)
