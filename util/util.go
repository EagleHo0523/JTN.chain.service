package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type dataBase struct {
	Url      string `json:"Url"`
	DBname   string `json:"DBname"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type ChainSetting struct {
	Env       string `json:"Env"`
	Url       string `json:"Url"`
	LoopDelay int
	db        dataBase
}
type MqttSetting struct {
	Env      string `json:"Env"`
	Url      string `json:"Url"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type envSetting struct {
	DDMX ChainSetting
	ETH  ChainSetting
	BTC  ChainSetting
	MQTT MqttSetting
}
type Addrs struct {
	Host string
	Addr []string `json:",omitempty"`
}
type BackTxMessage struct {
	To          string
	From        string
	BlockNumber string
	TxHash      string
	Symbol      string
	Amount      string
	Contract    string
}
type ListenTxReturn struct {
	HostName string
	TxBack   BackTxMessage
}

func FailOnError(err error, msg string) {
	if err != nil {
		// log.Fatalf("%s: %s", msg, err)
		fmt.Printf("%s: %s\n", msg, err)
	}
}
func GetSystemParameters(method string) (interface{}, error) {
	var rtn interface{}

	jsonFile, err := os.Open("./env.json")
	if err != nil {
		return rtn, err
	}
	byteVal, _ := ioutil.ReadAll(jsonFile)

	var envs envSetting
	err = json.Unmarshal(byteVal, &envs)
	if err != nil {
		return rtn, err
	}

	switch method {
	case "DDMX":
		rtn = envs.DDMX
	case "ETH":
		rtn = envs.ETH
	case "BTC":
		rtn = envs.BTC
	case "MQTT":
		rtn = envs.MQTT
	}

	return rtn, nil
}

// big.Int to Hex string
func BigInt2HexString(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %X or upper case
}
func RpcPost(url string, args string) (string, error) {
	if url == "" {
		err := errors.New("Empty URL")
		return "", err
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(args))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func Hex2AsciiString(s string) string {
	var rtn = ""
	bs, err := hex.DecodeString(s)
	if err != nil {
		return rtn
	}
	rtn = string(bs)
	return rtn
}
func EncryptKey(text string, passphrase string) (string, error) {
	data := []byte(text)
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	noneSize := gcm.NonceSize()
	fmt.Println("noneSize: ", noneSize)
	nonce := make([]byte, noneSize)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	seal := gcm.Seal(nonce, nonce, data, nil)

	return base58.Encode(seal), nil
}
func DecryptKey(text string, passphrase string) (string, error) {
	data := []byte(base58.Decode(text))
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// IsZeroAddress validate if it's a 0 address
func IsZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// CalcGasCost calculate gas cost given gas limit (units) and gas price (wei)
func CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	gasLimitBig := big.NewInt(int64(gasLimit))
	return gasLimitBig.Mul(gasLimitBig, gasPrice)
}
