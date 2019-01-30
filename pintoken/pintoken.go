package main
import (
  "encoding/base64"
	"crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "crypto/rsa"
  "crypto/sha256"
  "crypto/x509"
  "fmt"
  "encoding/pem"
  config "github.com/wenewzhang/mixin_labs-go-bot/config"
)
func main() {
  var token []byte
  var Msg error
  // var n1,n2 int
  // n1,n2 = 3,4
  fmt.Println(config.PinToken)
  token,Msg = base64.StdEncoding.DecodeString(config.PinToken)
  fmt.Println("token:",token)
  fmt.Println("token length:",len(token))
  if Msg != nil {
    fmt.Println("Pin token base64 format error!",Msg)
  }
  if len(token) != 40 {
		fmt.Println("Pin token base64 data must equal 40!")
  }
  privBlock, _ := pem.Decode([]byte(config.PrivateKey))
  if privBlock == nil {
    fmt.Println("invalid pem private key")
  }
  priv, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
  if err != nil {
    fmt.Println("ParsePKCS1PrivateKey error")
  }
  // token, _ := base64.StdEncoding.DecodeString(pinToken)
  keyBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, token, []byte(config.SessionId))
  if err != nil {
    fmt.Println("DecryptOAEP error")
  }
  // pinByte := []byte(pin)

  // key, err := base64.StdEncoding.DecodeString(current.PINSecret)
  // if err != nil {
  //   return "", session.ServerError(ctx, err)
  // }
  // if len(key) != 40 || len(pinBytes) < aes.BlockSize {
  //   return "", session.BadDataError(ctx)
  // }
  //PIN generate by python client,could replace by other
  PinData := "fFH5J54zjaMN09+s/b1+NlYtK2q9KKp9YH2y5YyvbfiWqhEP2cEBF65x6me6qcFq"
  pinBytes,err := base64.StdEncoding.DecodeString(PinData)
	if err != nil {
		fmt.Println("Pin data base64 data error!")
	}
  fmt.Println(pinBytes)
  fmt.Println(len(pinBytes))
  if len(pinBytes)%aes.BlockSize != 0 {
    fmt.Println("Pin Data must multiple 16!")
  }
  // fmt.Println("sub array of 0-31:",PinCode[:32])
  // block, err := aes.NewCipher(PinCode[:32])
  // if err != nil {
  //   fmt.Println("aes init error!")
  // }
  // fmt.Println("block:",block)

  block, err := aes.NewCipher(keyBytes[:32])
  if err != nil {
    fmt.Println("get 32bit PINSECRET and NEW a AES block error!",err)
  }

  iv := pinBytes[:aes.BlockSize]
  pin_timestamp := pinBytes[aes.BlockSize:]
  fmt.Println("iv 16 length is bit:",iv)
  fmt.Println("pin and timestamp:",pin_timestamp)
  mode := cipher.NewCBCDecrypter(block, iv)
  mode.CryptBlocks(pin_timestamp, pin_timestamp)
  length := len(pin_timestamp)
  unpadding := int(pin_timestamp[length-1])
  fmt.Println(unpadding)
  if unpadding > length {
    fmt.Println("unpadding must small than length of pin&timestamp")
  }
}
