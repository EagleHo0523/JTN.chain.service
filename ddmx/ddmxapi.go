package ddmx

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	mqtt "../mqtt"
	util "../util"

	m "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
)

var mqSetting util.MqttSetting
var mqttConnect m.Client

var listenTxBook util.ValueDictionary
var maxBlockNum int64
var nowBlockNum int64

var IsSetting bool = false
var url string = ""
var key string = ""
var delay time.Duration = time.Second * time.Duration(5)

type responseNewAccount struct {
	NewAccount string
}
type responseBalance struct {
	Balance string
	Symbol  string
}
type resopnseTxFee struct {
	Normal string
	High   string
}
type responseBlockNumber struct {
	BlockNumber string
}
type responseTransaction struct {
	TxHash      string
	BlockNumber string
	TxFee       string
}
type responseAddToken struct {
	Balance string
	Symbol  string
}
type responseSetListenTxAddress struct {
	Code    int64
	Message string
}
type responseGetListenTxAddress struct {
	Address []string
}

type requestNewAccount struct {
	Uuid     string
	Password string
}
type requestGetBalance struct {
	Uuid string
	Data dataGetBalance
}
type dataGetBalance struct {
	Account  string
	Contract string
}
type requestSendTransaction struct {
	Uuid string
	Data dataSendTransaction
}
type dataSendTransaction struct {
	From     string
	To       string
	Amount   string
	Password string
	Contract string
}
type requestAddToken struct {
	Uuid string
	Data dataAddToken
}
type dataAddToken struct {
	Account  string
	Contract string
}
type requestSetListenTxAddress struct {
	Host    string
	Address []string
}
type requestGetListenTxAddress struct {
	Host string
}

type paramsNewAccount struct {
	uuid     string
	password string
}
type getBalanceInfo struct {
	account  string
	contract string
}
type paramsGetBalance struct {
	uuid string
	more getBalanceInfo
}
type sendTransactionInfo struct {
	from     string
	to       string
	amount   string
	password string
	contract string
}
type paramsSendTransaction struct {
	uuid string
	more sendTransactionInfo
}

func ProcessMethod(method string, params interface{}) (interface{}, error) {
	var rtn interface{}
	var err error = nil

	if !IsSetting {
		return nil, errors.New("server error: loss environment settings")
	}

	start := time.Now()
	// method = strings.ToUpper(method)
	switch strings.ToUpper(method) {
	case "NEWACCOUNT":
		rtn, err = processNewAccount(params)
	case "GETBALANCE":
		rtn, err = processGetBalance(params)
	case "GETTXFEE":
		rtn, err = processGetTxFee(params)
	case "GETBLOCKNUMBER":
		rtn, err = processGetBlockNumber(params)
	case "SENDTRANSACTION":
		rtn, err = processSendTransaction(params)
	case "ADDTOKEN":
		rtn, err = processAddToken(params)
	case "SETLISTENTXADDRESS":
		rtn, err = processSetListenTxAddress(params)
	case "GETLISTENTXADDRESS":
		rtn, err = processGetListenTxAddress(params)
	}
	if err != nil {
		return rtn, err
	}
	log.Printf("ChainAPI: %s\t%s", method, time.Since(start))
	return rtn, nil
}
func ProcessTxBack() {
	p, err := util.GetSystemParameters("DDMX")
	if err != nil {
		util.FailOnError(err, "Setting DDMX environment parameters fail.")
		return
	} else {
		env := p.(util.ChainSetting)
		sInit(env.Url, env.LoopDelay)

		m, err := util.GetSystemParameters("MQTT")
		if err != nil {
			util.FailOnError(err, "Setting DDMX environment parameters fail.")
			return
		}
		mqSetting = m.(util.MqttSetting)
		mqttConnect, err = mqtt.CreateConnect(mqSetting)
		if err != nil {
			util.FailOnError(err, "create mqtt connection fail.")
			return
		}

		num, err := Rpc_blockNumber()
		if err != nil {
			return
		}
		maxBlockNum = num
		nowBlockNum = num

		fmt.Println("Start to sync DDMX chain transaction...")
		go getBlockNum()
		listenTxBack()
	}
}

func sInit(sUrl string, sDelay int) {
	if sUrl != "" && sDelay > 0 {
		url = sUrl
		delay = time.Second * time.Duration(sDelay)
		IsSetting = true
	}
}
func processNewAccount(params interface{}) (interface{}, error) {
	var respNewAcc responseNewAccount

	// 檢查型態
	tmp, ok := params.(map[string]interface{})
	if !ok {
		err := errors.New("Error parameter")
		return nil, err
	}

	var reqNewAccount requestNewAccount
	if err := mapstructure.Decode(tmp, &reqNewAccount); err != nil {
		err := errors.New("Error parameter")
		return nil, err
	}

	passphrase := reqNewAccount.Password

	account, err := Rpc_newAccount(passphrase)
	if err != nil {
		return nil, err
	}
	respNewAcc.NewAccount = account

	return respNewAcc, nil
}
func processGetBalance(params interface{}) (interface{}, error) {
	respBalance := responseBalance{Balance: "", Symbol: "DDMX"}

	// 檢查型態
	tmp, ok := params.(map[string]interface{})
	if !ok {
		err := errors.New("Error parameter")
		return nil, err
	}

	var reqGetBalance requestGetBalance
	if err := mapstructure.Decode(tmp, &reqGetBalance); err != nil {
		err := errors.New("Error parameter")
		return nil, err
	}

	account := reqGetBalance.Data.Account
	contract := reqGetBalance.Data.Contract

	if contract == "" {
		accInfo, err := Rpc_getAccountInfo(account)
		if err != nil {
			return respBalance, err
		}
		respBalance.Balance = fmt.Sprintf("%v", accInfo.Balance)
	} else {
		balance, err := Rpc_getBalanceToken(account, contract)
		if err != nil {
			return respBalance, err
		}
		token, err := Rpc_getContractToken(contract)
		if err != nil {
			return respBalance, err
		}

		respBalance.Balance = strconv.FormatFloat(balance, 'f', -1, 64)
		respBalance.Symbol = token
	}

	return respBalance, nil
}
func processGetTxFee(params interface{}) (interface{}, error) {
	var respTxFee resopnseTxFee

	gasprice, err := Rpc_gasPrice()
	if err != nil {
		return respTxFee, err
	}

	normal := util.CalcGasCost(uint64(30000), gasprice)
	high := util.CalcGasCost(uint64(60000), gasprice)

	respTxFee.Normal = util.ToDecimal(normal, 18).String()
	respTxFee.High = util.ToDecimal(high, 18).String()

	return respTxFee, nil
}
func processGetBlockNumber(params interface{}) (interface{}, error) {
	var respBlockNumber responseBlockNumber

	b, err := Rpc_getBlockByNumber("latest")
	if err != nil {
		return respBlockNumber, err
	}

	blockNumber := big.NewInt(0)
	blockNumber, _ = blockNumber.SetString(b.Number[2:], 16)
	respBlockNumber.BlockNumber = fmt.Sprintf("%v", blockNumber)

	return respBlockNumber, nil
}
func processSendTransaction(params interface{}) (interface{}, error) {
	var respTransaction responseTransaction

	// 檢查型態
	tmp, ok := params.(map[string]interface{})
	if !ok {
		err := errors.New("Error parameter")
		return nil, err
	}

	var reqSendTransaction requestSendTransaction
	if err := mapstructure.Decode(tmp, &reqSendTransaction); err != nil {
		err := errors.New("Error parameter")
		return nil, err
	}

	from := reqSendTransaction.Data.From
	to := reqSendTransaction.Data.To
	amount, _ := strconv.ParseFloat(reqSendTransaction.Data.Amount, 64)
	password := reqSendTransaction.Data.Password
	contract := reqSendTransaction.Data.Contract

	passphrase := password

	if contract == "" {
		s, err := Rpc_sendTxWithPassword(from, to, amount, passphrase)
		if err != nil {
			return respTransaction, err
		}
		respTransaction.TxHash = s
	} else {
		s, err := Rpc_sendTxWithPasswordToken(from, to, amount, passphrase, contract)
		if err != nil {
			return respTransaction, err
		}
		respTransaction.TxHash = s
	}

	var getTx bool = false
	count := 0
	for respTransaction.TxHash != "" && !getTx {
		if count > 60 {
			break
		}
		time.Sleep(2500 * time.Millisecond)
		txInfo, err := Rpc_getTxByHash(respTransaction.TxHash)
		if err != nil {
			return respTransaction, err
		}
		if txInfo.BlockNumber != "" {
			num, _ := strconv.ParseInt(txInfo.BlockNumber[2:], 16, 64)
			respTransaction.BlockNumber = strconv.FormatInt(num, 10)

			gas, _ := strconv.ParseInt(txInfo.Gas[2:], 16, 64)
			gasPrice, _ := strconv.ParseInt(txInfo.GasPrice[2:], 16, 64)

			gasWei := util.CalcGasCost(uint64(gas), big.NewInt(gasPrice))
			respTransaction.TxFee = util.ToDecimal(gasWei, 18).String()
			getTx = true
		}
		count++
	}

	return respTransaction, nil
}
func processAddToken(params interface{}) (interface{}, error) {
	var respAddToken responseAddToken

	// 檢查型態
	tmp, ok := params.(map[string]interface{})
	if !ok {
		err := errors.New("Error parameter")
		return nil, err
	}

	var reqAddToken requestAddToken
	if err := mapstructure.Decode(tmp, &reqAddToken); err != nil {
		err := errors.New("Error parameter")
		return nil, err
	}

	account := reqAddToken.Data.Account
	contract := reqAddToken.Data.Contract

	symbol, err := Rpc_getContractToken(contract)
	if err != nil {
		return respAddToken, err
	}
	respAddToken.Symbol = symbol

	balance, err := Rpc_getBalanceToken(account, contract)
	if err != nil {
		return respAddToken, err
	}

	respAddToken.Balance = strconv.FormatFloat(balance, 'f', -1, 64)

	return respAddToken, nil
}
func processSetListenTxAddress(params interface{}) (interface{}, error) {
	respSetListenTxAddress := responseSetListenTxAddress{Code: 500, Message: "Setting addresses fail."}

	// 檢查型態
	tmp, ok := params.(map[string]interface{})
	if !ok {
		err := errors.New("Error parameter")
		return nil, err
	}

	var reqSetListenTxAddress requestSetListenTxAddress
	if err := mapstructure.Decode(tmp, &reqSetListenTxAddress); err != nil {
		err := errors.New("Error parameter")
		return nil, err
	}

	host := reqSetListenTxAddress.Host
	addrs := reqSetListenTxAddress.Address

	var count int64 = 0
	for i := 0; i < len(addrs); i++ {
		addr := addrs[i]
		if !listenTxBook.Has(addr) {
			var val []string
			val = append(val, host)
			listenTxBook.Set(addr, val)
			count++
		} else {
			val := listenTxBook.Get(addr).([]string)
			var ok bool = false
			for j := 0; j < len(val); j++ {
				if val[j] == host {
					ok = true
				}
			}
			if !ok {
				val = append(val, host)
				listenTxBook.Delete(addr)
				listenTxBook.Set(addr, val)
				count++
			}
		}
	}

	respSetListenTxAddress.Code = 200
	respSetListenTxAddress.Message = "Success, add " + strconv.FormatInt(count, 10) + " new addresses."
	fmt.Println("ETH chain has " + strconv.Itoa(listenTxBook.Size()) + " addresses for sync transaction.")

	return respSetListenTxAddress, nil
}
func processGetListenTxAddress(params interface{}) (interface{}, error) {
	var respGetListenTxAddr responseGetListenTxAddress

	// 檢查型態
	tmp, ok := params.(map[string]interface{})
	if !ok {
		err := errors.New("Error parameter")
		return nil, err
	}

	var reqGetListenTxAddress requestGetListenTxAddress
	if err := mapstructure.Decode(tmp, &reqGetListenTxAddress); err != nil {
		err := errors.New("Error parameter")
		return nil, err
	}

	host := reqGetListenTxAddress.Host

	keys := listenTxBook.Keys()
	if len(keys) > 0 {
		for i := 0; i < len(keys); i++ {
			s := fmt.Sprintf("%v", keys[i])
			val := listenTxBook.Get(s).([]string)
			for j := 0; j < len(val); j++ {
				if val[j] == host {
					respGetListenTxAddr.Address = append(respGetListenTxAddr.Address, s)
				}
			}
		}
	}

	return respGetListenTxAddr, nil
}

func listenTxBack() {
	for {
		if nowBlockNum > maxBlockNum {
			time.Sleep(delay)
			continue
		}

		numHex := "0x" + fmt.Sprintf("%x", nowBlockNum)
		blockInfo, err := Rpc_getBlockByNumber(numHex)
		if err != nil {
			util.FailOnError(err, "Get DDMX block info fail.")
			continue
		}

		addrs := listenTxBook.Keys()
		tsc := len(blockInfo.Transactions)
		asc := len(addrs)
		var txReturns []util.ListenTxReturn

		if tsc > 0 { //有交易紀錄
			for i := 0; i < tsc; i++ {
				to := blockInfo.Transactions[i].To
				if len(blockInfo.Transactions[i].Input) > 2 && strings.Index(blockInfo.Transactions[i].Input, "0xa9059cbb") >= 0 { // token 交易
					to, _, _ = analyzeInput(blockInfo.Transactions[i].Input)
				}
				for j := 0; j < asc; j++ {
					if strings.ToUpper(to) == strings.ToUpper(addrs[j].(string)) {
						hosts := listenTxBook.Get(addrs[j]).([]string)
						util.FailOnError(err, "Translate DDMX block number fail.")
						for k := 0; k < len(hosts); k++ {
							var txRtn util.ListenTxReturn

							txBack, err := toListenTx(blockInfo.Transactions[i])
							util.FailOnError(err, "Process DDMX listen tx infomation fail.")

							txRtn.HostName = hosts[k]
							txRtn.TxBack = txBack

							txReturns = append(txReturns, txRtn)
						}
					}
				}
			}
		} else {

		}

		if len(txReturns) > 0 {
			sendListenTx(txReturns)
		}

		nowBlockNum++
	}
}
func getBlockNum() {
	for {
		time.Sleep(1 * time.Second)
		num, err := Rpc_blockNumber()
		if err != nil {
			continue
		} else {
			if num > maxBlockNum {
				snedListenBlockNum(num)
			}
			maxBlockNum = num
		}
	}
}
func sendListenTx(arrs []util.ListenTxReturn) {
	if mqttConnect.IsConnectionOpen() {
		c := len(arrs)
		for i := 0; i < c; i++ {
			topic := "blockchain/ddmx/listentx/" + arrs[i].HostName

			body, err := json.Marshal(arrs[i].TxBack)
			util.FailOnError(err, "Failed to marshal JSON")

			go mqtt.Send(mqttConnect, topic, body)
		}
	} else {
		err := errors.New("connection not ready.")
		util.FailOnError(err, "mqtt connection fail.")
	}
}
func snedListenBlockNum(num int64) {
	if mqttConnect.IsConnectionOpen() {
		topic := "blockchain/ddmx/listenblocknumber"

		var blockNum responseBlockNumber
		blockNum.BlockNumber = fmt.Sprintf("%v", num)

		body, err := json.Marshal(blockNum)
		util.FailOnError(err, "Failed to marshal JSON")

		go mqtt.Send(mqttConnect, topic, body)
	} else {
		err := errors.New("connection not ready.")
		util.FailOnError(err, "mqtt connection fail.")
	}
}
func toListenTx(txInfo DDM_TransactionInfo) (util.BackTxMessage, error) {
	var backTxMessage util.BackTxMessage

	if strings.Index(txInfo.Input, "0xa9059cbb") >= 0 {
		// ERC20 token
		to, amount, err := analyzeInput(txInfo.Input)
		util.FailOnError(err, "Analysis the token information fail.")

		backTxMessage.Contract = txInfo.To
		backTxMessage.To = to
		backTxMessage.Amount = strconv.FormatFloat(amount, 'f', -1, 64)

		token, err := Rpc_getContractToken(backTxMessage.Contract)
		if err != nil {
			return backTxMessage, err
		}
		backTxMessage.Symbol = token
	} else {
		backTxMessage.To = txInfo.To

		s := txInfo.Value[2:]
		bInt := big.NewInt(0)
		bInt, _ = bInt.SetString(s, 16) // hex string to big.Int
		val := util.ToDecimal(bInt, 18)
		v, _ := val.Float64()
		backTxMessage.Amount = strconv.FormatFloat(v, 'f', -1, 64)

		backTxMessage.Symbol = "DDMX"
	}

	blockNumber := big.NewInt(0)
	blockNumber, _ = blockNumber.SetString(txInfo.BlockNumber[2:], 16)
	backTxMessage.BlockNumber = fmt.Sprintf("%v", blockNumber)
	backTxMessage.From = txInfo.From
	backTxMessage.TxHash = txInfo.Hash

	return backTxMessage, nil
}
func analyzeInput(input string) (string, float64, error) {
	var to string = "0x" + input[34:74]
	var value float64 = 0

	s := input[74:]
	bInt := big.NewInt(0)
	bInt, _ = bInt.SetString(s, 16) // hex string to big.Int
	val := util.ToDecimal(bInt, 18)
	value, _ = val.Float64()

	return to, value, nil
}

// func ddmHashPassphrase(s string) string {
// 	// sha256
// 	s = key + "+" + s + "+jutainet"
// 	h := sha256.New()
// 	h.Write([]byte(s))
// 	bs := h.Sum(nil)
// 	return fmt.Sprintf("%x", bs)
// }
