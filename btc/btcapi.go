package btc

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	util "../util"
)

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

var mqSetting util.MqttSetting

var IsSetting bool = false
var url string = ""
var key string = ""
var delay time.Duration = time.Second * time.Duration(5)

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
	}
	if err != nil {
		return rtn, err
	}
	log.Printf("ChainAPI: %s\t%s", method, time.Since(start))
	return rtn, nil
}
func ProcessTxBack() {
	p, err := util.GetSystemParameters("BTC")
	if err != nil {
		util.FailOnError(err, "Setting BTC environment parameters fail.")
		return
	} else {
		env := p.(util.ChainSetting)
		sInit(env.Url, env.LoopDelay)
	}

	go getBlockNum()
	listenTxBack()
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

	arr, ok := params.([]interface{})
	if !ok {
		err := errors.New("{\"Error parameter\"}")
		return nil, err
	}
	if len(arr) != 2 {
		err := errors.New("{\"Error parameter\"}")
		return nil, err
	}

	walletname := arr[0].(string) // uuid for walletname
	password := arr[1].(string)

	unloadAllWallet() // unload 所有錢包
	passpharse := password

	ok, err := Rpc_createwallet(walletname)
	if !ok {
		unloadAllWallet() // unload 所有錢包
		return respNewAcc, err
	}
	ok, err = Rpc_encryptwallet(passpharse)
	if !ok {
		unloadAllWallet() // unload 所有錢包
		return respNewAcc, err
	}
	account, err := Rpc_getnewaddress()
	if err != nil {
		unloadAllWallet() // unload 所有錢包
		return respNewAcc, err
	}

	respNewAcc.NewAccount = account

	unloadAllWallet() // unload 所有錢包

	return respNewAcc, nil
}
func processGetBalance(params interface{}) (interface{}, error) {
	respBalance := responseBalance{Balance: "", Symbol: "BTC"}

	arr, ok := params.([]interface{})
	if !ok {
		err := errors.New("{\"Error parameter\"}")
		return nil, err
	}
	if len(arr) != 2 {
		err := errors.New("{\"Error parameter\"}")
		return nil, err
	}

	walletname := arr[0].(string) // uuid for walletname

	unloadAllWallet() // unload 所有錢包

	ok, err := Rpc_loadwallet(walletname)
	if !ok {
		return respBalance, err
	}
	balance, err := Rpc_getbalance()
	if err != nil {
		return respBalance, err
	}

	respBalance.Balance = strconv.FormatFloat(balance, 'f', -1, 64)

	unloadAllWallet() // unload 所有錢包

	return respBalance, nil
}
func processGetTxFee(params interface{}) (interface{}, error) {
	var respTxFee resopnseTxFee
	return respTxFee, nil
}
func processGetBlockNumber(params interface{}) (interface{}, error) {
	var respBlockNumber responseBlockNumber

	b, err := Rpc_getblockcount()
	if err != nil {
		return respBlockNumber, err
	}

	blockNumber := big.NewInt(b)
	respBlockNumber.BlockNumber = fmt.Sprintf("%v", blockNumber)

	return respBlockNumber, nil
}
func processSendTransaction(params interface{}) (interface{}, error) {
	var rtn responseTransaction

	arr, ok := params.([]interface{})
	if !ok {
		err := errors.New("{\"Error parameter\"}")
		return nil, err
	}
	if len(arr) != 2 {
		err := errors.New("{\"Empty parameter\"}")
		return rtn, err
	}

	walletname := arr[0].(string) // uuid for walletname
	if _, ok := arr[1].(map[string]interface{}); !ok {
		err := errors.New("{\"Error parameter\"}")
		return rtn, err
	}
	tmp := arr[1].(map[string]interface{})
	from := tmp["FROM"].(string)
	to := tmp["TO"].(string)
	amount, _ := strconv.ParseFloat(tmp["AMOUNT"].(string), 64)
	password := tmp["PASSWORD"].(string)
	// contract := tmp["CONTRACT"].(string) // BTC 不提供

	unloadAllWallet() // unload 所有錢包

	_, err := Rpc_loadwallet(walletname)
	if err != nil {
		unloadAllWallet() // unload 所有錢包
		return rtn, err
	}

	_, err = checkAccountInWallet(walletname, from)
	if err != nil {
		err := errors.New("{\"from account not in wallet.\"}")
		unloadAllWallet() // unload 所有錢包
		return rtn, err
	}
	passphrase := password

	preTx, err := Rpc_createrawtransaction(to, amount)
	if err != nil {
		unloadAllWallet() // unload 所有錢包
		return rtn, err
	}
	fundRawTx, err := Rpc_fundrawtransaction(preTx)
	if err != nil {
		unloadAllWallet() // unload 所有錢包
		return rtn, err
	}
	fundHash := fundRawTx.Hex
	ok, err = Rpc_walletpassphrase(passphrase)
	if !ok {
		unloadAllWallet() // unload 所有錢包
		return rtn, err
	}
	signRawTx, err := Rpc_signrawtransactionwithwallet(fundHash)
	if !ok {
		unloadAllWallet() // unload 所有錢包
		return rtn, err
	}
	ok, err = Rpc_walletlock()
	if !ok {
		unloadAllWallet() // unload 所有錢包
		return rtn, err
	}
	txHash, err := Rpc_sendrawtransaction(signRawTx)
	if err != nil {
		unloadAllWallet() // unload 所有錢包
		return rtn, err
	}

	var getTx bool = false
	for rtn.BlockNumber == "" && !getTx {
		time.Sleep(2500 * time.Second)
		txInfo, err := Rpc_gettransaction(txHash)
		if err != nil {
			unloadAllWallet() // unload 所有錢包
			return rtn, err
		}
		if txInfo.Blockhash != "" {
			blockInfo, err := Rpc_getblock(txInfo.Blockhash)
			if err != nil {
				unloadAllWallet() // unload 所有錢包
				return rtn, err
			}
			rtn.BlockNumber = strconv.FormatInt(blockInfo.Height, 10)
			rtn.TxFee = strconv.FormatFloat(fundRawTx.Fee, 'f', -1, 64)
			rtn.TxHash = txInfo.Txid

			getTx = true
		}
	}

	unloadAllWallet() // unload 所有錢包

	return rtn, nil
}
func processAddToken(params interface{}) (interface{}, error) {
	var respAddToken responseAddToken
	return respAddToken, nil
}

func listenTxBack() {

}
func getBlockNum() {

}
func checkAccountInWallet(walletname string, account string) (bool, error) {
	rtn := false

	ok, err := Rpc_loadwallet(walletname)
	if !ok {
		return rtn, err
	}
	accounts, err := Rpc_getaddressesbylabel("")
	for i := range accounts {
		if strings.Index(accounts[i], account) >= 0 {
			rtn = true
		}
	}

	return rtn, nil
}
func unloadAllWallet() {
	wallets, err := Rpc_listwallets()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		for i := 0; i < len(wallets); i++ {
			name := wallets[i]
			Rpc_unloadwallet(name)
		}
	}
}

// func btcHashPassphrase(s string) string {
// 	// sha256
// 	s = key + "+" + s + "+jutainet"
// 	h := sha256.New()
// 	h.Write([]byte(s))
// 	bs := h.Sum(nil)
// 	return fmt.Sprintf("%x", bs)
// }
