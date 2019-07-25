package btc

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	util "../util"
)

type BTC_FundRawTransaction struct {
	Hex       string
	Fee       float64
	Changepos float64
}
type BTC_BlockInfo struct {
	Hash              string
	Confirmations     int64
	Strippedsize      int64
	Size              int64
	Weight            int64
	Height            int64
	Version           int64
	Versionhex        string
	Merklerroot       string
	Tx                []string `json:",omitempty"`
	Time              int64
	Mediantime        int64
	Nonce             int64
	Bits              string
	Difficulty        float64
	Chainwork         string
	Ntx               int64
	Previousblockhash string
}
type BTC_TransactionInfo struct {
	Amount             float64
	Fee                float64
	Confirmations      int64
	Blockhash          string
	Blockindex         int64
	Blocktime          int64
	Trusted            bool
	Txid               string
	Walletconflicts    interface{} `json:",omitempty"`
	Time               int64
	Timereceived       int64
	Bip125_replaceable string
	Details            []BTC_TransactionDetail `json:",omitempty"`
	Hex                string
}
type BTC_TransactionDetail struct {
	Address   string
	Category  string
	Amount    float64
	Label     string
	Vout      int64
	Fee       float64
	Abandoned bool
}

type btcError struct {
	Code    int64
	Message string
}
type btcSignTxInfo struct {
	Hex      string
	Complete bool
}
type btcCreateWalletInfo struct {
	Name    string
	Warning string
}
type btcBalanceResp struct {
	Result float64  `json:",omitempty"`
	Error  btcError `json:",omitempty"`
	Id     string
}
type btcCreateRawTxResp struct {
	Result string   `json:",omitempty"`
	Error  btcError `json:",omitempty"`
	Id     string
}
type btcFundRawTxResp struct {
	Result BTC_FundRawTransaction `json:",omitempty"`
	Error  btcError               `json:",omitempty"`
	Id     string
}
type btcSignRawTxWithWalletResp struct {
	Result btcSignTxInfo `json:",omitempty"`
	Error  btcError      `json:",omitempty"`
	Id     string
}
type btcSendRawTxResp struct {
	Result string   `json:",omitempty"`
	Error  btcError `json:",omitempty"`
	Id     string
}
type btcWalletPassphraseResp struct {
	Result interface{} `json:",omitempty"`
	Error  btcError    `json:",omitempty"`
	Id     string
}
type btcWalletLockResp struct {
	Result interface{} `json:",omitempty"`
	Error  btcError    `json:",omitempty"`
	Id     string
}
type btcListWalletsResp struct {
	Result []string `json:",omitempty"`
	Error  btcError `json:",omitempty"`
	Id     string
}
type btcGetAddressesByLabelResp struct {
	Result interface{} `json:",omitempty"`
	Error  btcError    `json:",omitempty"`
	Id     string
}
type btcBlockCountResp struct {
	Result int64    `json:",omitempty"`
	Error  btcError `json:",omitempty"`
	Id     string
}
type btcBlockHashResp struct {
	Result string   `json:",omitempty"`
	Error  btcError `json:",omitempty"`
	Id     string
}
type btcBlockInfoResp struct {
	Result BTC_BlockInfo `json:",omitempty"`
	Error  btcError      `json:",omitempty"`
	Id     string
}
type btcTxInfoResp struct {
	Result BTC_TransactionInfo `json:",omitempty"`
	Error  btcError            `json:",omitempty"`
	id     string
}
type btcCreateWalletResp struct {
	Result btcCreateWalletInfo `json:",omitempty"`
	Error  btcError            `json:",omitempty"`
	Id     string
}
type btcLoadWalletResp struct {
	Result interface{} `json:",omitempty"`
	Error  btcError    `json:",omitempty"`
	Id     string
}
type btcUnloadWalletResp struct {
	Result interface{} `json:",omitempty"`
	Error  btcError    `json:",omitempty"`
	Id     string
}
type btcEncryptWalletResp struct {
	Result string   `json:",omitempty"`
	Error  btcError `json:",omitempty"`
	Id     string
}
type btcWalletPassphraseChangeResp struct {
	Result interface{} `json:",omitempty"`
	Error  btcError    `json:",omitempty"`
	Id     string
}
type btcGetNewAddressResp struct {
	Result string   `json:",omitempty"`
	Error  btcError `json:",omitempty"`
	Id     string
}

func processError(m interface{}) error {
	e := m.(btcError)
	code := fmt.Sprintf("%v", e.Code)
	message := e.Message
	err := errors.New("{\"CODE\":\"" + code + "\",\"MESSAGE\":\"" + message + "\"}")
	return err
}

func Rpc_loadwallet(walletname string) (bool, error) {
	rtn := false

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"loadwallet\", \"params\":[\"" + walletname + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcLoadWalletResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		s := bResp.Error.Message
		if strings.Index(s, "Wallet file verification failed: Error loading wallet "+walletname+". Duplicate -wallet filename specified.") >= 0 {
			return true, nil
		}
		return rtn, processError(bResp.Error)
	} else {
		rtn = true
	}

	return rtn, nil
}
func Rpc_unloadwallet(walletname string) (bool, error) {
	rtn := false

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"unloadwallet\", \"params\":[\"" + walletname + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcUnloadWalletResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		s := bResp.Error.Message
		if strings.Index(s, "Requested wallet does not exist or is not loaded") >= 0 {
			return true, nil
		}
		return rtn, processError(bResp.Error)
	} else {
		rtn = true

	}

	return rtn, nil
}
func Rpc_createwallet(walletname string) (bool, error) {
	rtn := false

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"createwallet\", \"params\":[\"" + walletname + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcCreateWalletResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		if bResp.Result.Name == walletname {
			rtn = true
		}
	}

	return rtn, nil
}
func Rpc_encryptwallet(passphrase string) (bool, error) {
	rtn := false

	args := "{\"jsonrpc\":\"1.0\",\"id\":\"jutainet\",\"method\":\"encryptwallet\",\"params\":[\"" + passphrase + "\"]}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcEncryptWalletResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		s := bResp.Result
		if strings.Index(s, "wallet encrypted") >= 0 {
			return true, nil
		}
	}

	return rtn, nil
}
func Rpc_walletpassphrasechange(oldpassphrase string, newpassphrase string) (bool, error) {
	rtn := false

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"walletpassphrasechange\", \"params\":[\"" + oldpassphrase + "\",\"" + newpassphrase + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcWalletPassphraseChangeResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = true
	}

	return rtn, nil
}
func Rpc_getnewaddress() (string, error) {
	rtn := ""

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"getnewaddress\", \"params\":[], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcGetNewAddressResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result
	}

	return rtn, nil
}
func Rpc_getbalance() (float64, error) {
	var rtn float64

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"getbalance\", \"params\":[\"*\",6], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcBalanceResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result
	}

	return rtn, nil
}
func Rpc_listwallets() ([]string, error) {
	var rtn []string

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"listwallets\", \"params\":[], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcListWalletsResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		for i := 0; i < len(bResp.Result); i++ {
			rtn = append(rtn, bResp.Result[i])
		}
	}

	return rtn, nil
}
func Rpc_getaddressesbylabel(label string) ([]string, error) {
	var rtn []string

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"getaddressesbylabel\", \"params\":[\"" + label + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcGetAddressesByLabelResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		if bResp.Result == nil {
			return rtn, nil
		} else {
			arr := bResp.Result.(map[string]interface{})
			for key, _ := range arr {
				rtn = append(rtn, key)
			}
		}
	}

	return rtn, nil
}
func Rpc_walletpassphrase(passphrase string) (bool, error) {
	// 解鎖錢包
	rtn := false

	args := "{\"jsonrpc\":\"1.0\",\"id\":\"jutainet\",\"method\":\"walletpassphrase\",\"params\":[\"" + passphrase + "\", 60]}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcWalletPassphraseResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = true
	}

	return rtn, nil
}
func Rpc_walletlock() (bool, error) {
	// 錢包上鎖
	rtn := false

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"walletlock\", \"params\":[], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcWalletLockResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = true
	}

	return rtn, nil
}
func Rpc_signrawtransactionwithwallet(fundTxHash string) (string, error) {
	rtn := ""

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"signrawtransactionwithwallet\", \"params\":[\"" + fundTxHash + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcSignRawTxWithWalletResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result.Hex
	}

	return rtn, nil
}
func Rpc_createrawtransaction(to string, value float64) (string, error) {
	rtn := ""

	val := strconv.FormatFloat(value, 'f', -1, 64)
	args := "{\"jsonrpc\":\"1.0\",\"id\":\"jutainet\",\"method\":\"createrawtransaction\",\"params\":[[],[{\"" + to + "\":" + val + "}]]}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcCreateRawTxResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result
	}

	return rtn, nil
}
func Rpc_fundrawtransaction(preTxHash string) (BTC_FundRawTransaction, error) {
	var rtn BTC_FundRawTransaction

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"fundrawtransaction\", \"params\":[\"" + preTxHash + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcFundRawTxResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result
	}

	return rtn, nil
}
func Rpc_sendrawtransaction(signTxHash string) (string, error) {
	rtn := ""

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"sendrawtransaction\", \"params\":[\"" + signTxHash + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcSendRawTxResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result
	}

	return rtn, nil
}
func Rpc_getblock(blockHash string) (BTC_BlockInfo, error) {
	var rtn BTC_BlockInfo

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"getblock\", \"params\":[\"" + blockHash + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcBlockInfoResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result
	}

	return rtn, nil
}
func Rpc_getblockcount() (int64, error) {
	var rtn int64

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"getblockcount\", \"params\":[], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcBlockCountResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result
	}

	return rtn, nil
}
func Rpc_getblockhash(blockHeight int64) (string, error) {
	rtn := ""

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"getblockhash\", \"params\":[" + strconv.FormatInt(blockHeight, 10) + "], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcBlockHashResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result
	}

	return rtn, nil
}
func Rpc_gettransaction(txHash string) (BTC_TransactionInfo, error) {
	var rtn BTC_TransactionInfo

	args := "{\"jsonrpc\":\"1.0\", \"method\":\"gettransaction\", \"params\":[\"" + txHash + "\"], \"id\":\"jutainet\"}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp btcTxInfoResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn = bResp.Result
	}

	return rtn, nil
}
