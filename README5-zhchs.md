# 通过 Go 买卖Bitcoin
![cover](https://github.com/wenewzhang/mixin_labs-go-bot/raw/master/Bitcoin_go.jpg)

上一章介绍了[Exincore](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README4-zhchs.md)，你可以1秒完成资产的市价买卖。如果你想限定价格买卖，或者买卖一些exincore不支持的资产，你需要OceanOne。

## 方案二: 挂单Ocean.One交易所
[Ocean.one](https://github.com/mixinNetwork/ocean.one)是基于Mixin Network的去中心化交易所，它性能一流。
你可以在OceanOne上交易任何资产，只需要将你的币转给OceanOne, 将交易信息写在交易的memo里，OceanOne会在市场里列出你的交易需求，
交易成功后，会将目标币转入到你的MixinNetwork帐上，它有三大特点与优势：
- 不需要在OceanOne注册
- 不需要存币到交易所
- 支持所有Mixin Network上能够转账的资产，所有的ERC20 EOS代币。

### 预备知识:
你先需要创建一个机器人, 方法在 [教程一](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README-zhchs.md).


#### 安装依赖包
[第四课](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/README4-zhchs.md), 在上一课中已经安装好了.

#### 充币到 Mixin Network, 并读出它的余额.
此处演示用 USDT购买BTC 或者 用BTC购买USDT。交易前，先检查一下钱包地址。
完整的步骤如下:
- 检查比特币或USDT的余额，钱包地址。并记下钱包地址。
- 从第三方交易所或者你的冷钱包中，将币充到上述钱包地址。
- 再检查一下币的余额，看到帐与否。(比特币的到帐时间是5个区块的高度，约100分钟)。

比特币与USDT的充值地址是一样的。

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

#### 取得Ocean.one的市场价格信息
如何来查询Ocean.one市场的价格信息呢？你要先了解你交易的基础币是什么，如果你想买比特币，卖出USDT,那么基础货币就是USDT;如果你想买USDT,卖出比特币，那么基础货币就是比特币.

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

#### 交易前，创建一个Memo!
在第二章里,[基于Mixin Network的 Java 比特币开发教程: 机器人接受比特币并立即退还用户](https://github.com/wenewzhang/mixin_labs-java-bot/blob/master/README2-zhchs.md), 我们学习过转帐，这儿我们介绍如何告诉Ocean.one，我们给它转帐的目的是什么，信息全部放在memo里.
- **Side** 方向,"B" 或者 "A", "B"是购买, "A"是出售.
- **AssetUUID** 目标虚拟资产的UUID.
- **Price** 价格，如果操作方向是"B", 价格就是AssetUUID的价格; 如果操作方向是"B", 价格就是转给Ocean.one币的价格.

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

#### 出售XIN的例子
转打算出售的 **XIN** 给Ocean.one(OCEANONE_BOT),将你打算换回来的目标虚拟资产的UUID放入memo.

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

如果你是打算买XIN,操作如下:
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

一个成功的挂单如下：

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

#### 取消挂单
Ocean.one将trace_id当做订单，比如上面的例子， **26ef0ec1-a60c-4702-b8ca-4bd40330e120** 就是订单号，我们用他来取消订单。
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
#### 通过读取资产余额，来确认到帐情况

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

## 源代码执行
编译执行，即可开始交易了.

- **go run coin_exchange.go**   编译项目.

本代码执行时的命令列表:

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

[完整代码](https://github.com/wenewzhang/mixin_labs-go-bot/blob/master/coin_exchange/coin_exchange.go)
