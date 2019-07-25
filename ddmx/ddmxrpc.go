package ddmx

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	util "../util"
)

type DDM_AccountInfo struct {
	Balance  float64
	Contract string
	Nonce    string
}
type DDM_BlockInfo struct {
	Difficulty          string
	ExtraData           string
	GasLimit            string
	GasUsed             string
	Hash                string
	Miner               string
	Number              string
	ParentHash          string
	Size                string
	StateRoot           string
	Timestamp           string
	TotalDifficulty     string
	TransactionReceipts []DDM_TransactionReceipts `json:",omitempty"`
	Transactions        []DDM_TransactionInfo     `json:",omitempty"`
	TransactionsRoot    string
}
type DDM_TransactionReceipts struct {
	Contractaddress string `json:",omitempty"`
	Gasused         string
	Logs            []DDM_ReceiptLogs `json:",omitempty"`
	Status          string
	Transactionhash string
}
type DDM_ReceiptLogs struct {
	Address          string
	Topics           []string `json:",omitempty"`
	Data             string
	BlockNumber      string
	TransactionHash  string
	TransactionIndex string
	BlockHash        string
	LogIndex         string
	Removed          bool
}
type DDM_TransactionInfo struct {
	BlockHash        string
	BlockNumber      string
	From             string
	Gas              string
	GasPrice         string
	Hash             string
	Input            string
	Nonce            string
	To               string
	TransactionIndex string
	Value            string
	V                string
	R                string
	S                string
}
type DDM_TransactionReceiptInfo struct {
	BlockHash         string
	BlockNumber       string
	ContractAddress   string `json:",omitempty"`
	CumulativeGasUsed string
	From              string
	GasUsed           string
	Logs              []string `json:",omitempty"`
	LogsBloom         string
	Status            string
	To                string
	TransactionHash   string
	TransactionIndex  string
}
type DDM_SignTxInfo struct {
	Raw string
	Tx  DDM_TxInfo
}
type DDM_TxInfo struct {
	Gas      string
	GasPrice string
	Hash     string
	Input    string
	Nonce    string
	To       string
	Value    string
	V        string
	R        string
	S        string
}

type ddmAccountInfo struct {
	Balance  big.Int
	Contract string
	Nonce    string
}
type ddmError struct {
	Code    int64
	Message string
}
type ddmAccounts struct {
	Jsonrpc string
	Id      int
	Result  []string `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmGetAccountInfoResp struct {
	Jsonrpc string
	Id      int
	Result  ddmAccountInfo `json:",omitempty"`
	Error   ddmError       `json:",omitempty"`
}
type ddmGetBalanceTokenResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmNewAccoutResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmGetBlockByNumberResp struct {
	Jsonrpc string
	Id      int
	Result  DDM_BlockInfo `json:",omitempty"`
	Error   ddmError      `json:",omitempty"`
}
type ddmGetBlockByHashResp struct {
	Jsonrpc string
	Id      int
	Result  DDM_BlockInfo `json:",omitempty"`
	Error   ddmError      `json:",omitempty"`
}
type ddmBlockTxCount struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmGetTxByHashResp struct {
	Jsonrpc string
	Id      int
	Result  DDM_TransactionInfo `json:",omitempty"`
	Error   ddmError            `json:",omitempty"`
}
type ddmGetTxReceiptByHashResp struct {
	Jsonrpc string
	Id      int
	Result  DDM_TransactionReceiptInfo `json:",omitempty"`
	Error   ddmError                   `json:",omitempty"`
}
type ddmSendTxResponse struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmSignTxWithPasswordResp struct {
	Jsonrpc string
	Id      int
	Result  DDM_SignTxInfo `json:",omitempty"`
	Error   ddmError       `json:",omitempty"`
}
type ddmSendRawTxResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmGasPriceResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmAccountsResp struct {
	Jsonrpc string
	Id      int
	Result  []string `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmNewAccountResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmGetBlockTxCountByNumberResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmGetBlockTxCountByHashResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmSendTxWithPasswordResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmSendTxWithPasswordTokenResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmGetPriceResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmBlockNumberResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmSendTransactionTokenResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmGetContractTokenResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}
type ddmGetContractNameResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ddmError `json:",omitempty"`
}

// return value = gas, gasPrice, nonce
func getSignTxInfo(from string) (string, string, string, error) {
	var gas string = "0xEA60" // 60000

	accountInfo, err := Rpc_getAccountInfo(from)
	if err != nil {
		return "", "", "", err
	}
	nonce := accountInfo.Nonce

	price, err := Rpc_gasPrice()
	if err != nil {
		return "", "", "", err
	}
	gasPrice := "0x" + util.BigInt2HexString(price)

	return gas, gasPrice, nonce, nil
}
func processError(m interface{}) error {
	e := m.(ddmError)
	code := fmt.Sprintf("%v", e.Code)
	message := e.Message
	err := errors.New("{\"CODE\":\"" + code + "\",\"MESSAGE\":\"" + message + "\"}")
	return err
}

func Rpc_accounts() ([]string, error) {
	var rtn []string

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_accounts\", \"PARAMS\":[], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmAccountsResp
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
func Rpc_getAccountInfo(account string) (DDM_AccountInfo, error) {
	var rtn DDM_AccountInfo

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_getAccountInfo\", \"PARAMS\":[\"" + account + "\"], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetAccountInfoResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		rtn.Contract = bResp.Result.Contract
		rtn.Nonce = bResp.Result.Nonce
		val := util.ToDecimal(&bResp.Result.Balance, 18)
		rtn.Balance, _ = val.Float64()
	}

	return rtn, nil
}
func Rpc_getBalanceToken(account string, contract string) (float64, error) {
	var rtn float64

	params := "[{\"DATA\":\"0x70a08231000000000000000000000000" + account[2:] + "\", \"TO\":\"" + contract + "\"},\"latest\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_call\", \"PARAMS\":" + params + ", \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetBalanceTokenResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		s := bResp.Result[2:]
		bInt := big.NewInt(0)
		bInt, _ = bInt.SetString(s, 16) // hex string to big.Int

		val := util.ToDecimal(bInt, 18)
		rtn, _ = val.Float64()
	}

	return rtn, nil
}
func Rpc_newAccount(passpharse string) (string, error) {
	var rtn string

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_newAccount\", \"PARAMS\":[\"" + passpharse + "\"], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmNewAccountResp
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
func Rpc_getBlockByNumber(num string) (DDM_BlockInfo, error) {
	var rtn DDM_BlockInfo

	// 分析 num 字串, 是否為數字; 如為文字, 是否為特定字串

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_getBlockByNumber\", \"PARAMS\":[\"" + num + "\",true], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetBlockByNumberResp
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
func Rpc_getBlockByHash(hashValue string) (DDM_BlockInfo, error) {
	var rtn DDM_BlockInfo

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_getBlockByHash\", \"PARAMS\":[\"" + hashValue + "\",true], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetBlockByHashResp
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
func Rpc_getBlockTxCountByNumber(num int64) (int64, error) {
	var rtn int64

	hexStr := strconv.FormatInt(num, 16) //int64 to hex string

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_getBlockTxCountByNumber\", \"PARAMS\":[\"0x" + hexStr + "\"], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetBlockTxCountByNumberResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		val := bResp.Result[2:]
		rtn, _ = strconv.ParseInt(val, 16, 64) // hex string to int64
	}

	return rtn, nil
}
func Rpc_getBlockTxCountByHash(hashValue string) (int64, error) {
	var rtn int64

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_getBlockTxCountByHash\", \"PARAMS\":[\"" + hashValue + "\"], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetBlockTxCountByHashResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		val := bResp.Result[2:]
		rtn, _ = strconv.ParseInt(val, 16, 64)
	}

	return rtn, nil
}
func Rpc_getTxByHash(hashValue string) (DDM_TransactionInfo, error) {
	var rtn DDM_TransactionInfo

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_getTxByHash\", \"PARAMS\":[\"" + hashValue + "\"], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetTxByHashResp
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
func Rpc_getTxReceiptByHash(hashValue string) (DDM_TransactionReceiptInfo, error) {
	var rtn DDM_TransactionReceiptInfo

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_getTxReceiptByHash\", \"PARAMS\":[\"" + hashValue + "\"], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetTxReceiptByHashResp
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
func Rpc_sendTxWithPassword(from string, to string, value float64, password string) (string, error) {
	var rtn string

	if len(from) != 42 {
		err := errors.New("{\"from: non-regular address.\"}")
		return "", err
	}
	if len(to) != 42 {
		err := errors.New("{\"to: non-regular address.\"}")
		return "", err
	}

	val := util.ToWei(value, 18)
	v := util.BigInt2HexString(val)

	params := "[{\"FROM\":\"" + from + "\", \"TO\":\"" + to + "\", \"VALUE\":\"0x" + v + "\"},\"" + password + "\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_sendTxWithPassword\", \"PARAMS\":" + params + ", \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmSendTxWithPasswordResp
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
func Rpc_sendTxWithPasswordToken(from string, to string, value float64, password string, contract string) (string, error) {
	var rtn string = ""

	if len(from) != 42 {
		err := errors.New("{\"from: non-regular address.\"}")
		return "", err
	}
	if len(to) != 42 {
		err := errors.New("{\"to: non-regular address.\"}")
		return "", err
	}
	if len(contract) != 42 {
		err := errors.New("{\"contract: non-regular address.\"}")
		return "", err
	}

	val := util.ToWei(value, 18)
	v := util.BigInt2HexString(val)
	amount := "0000000000000000000000000000000000000000000000000000000000000000" + v

	data := "0xa9059cbb000000000000000000000000" + to[2:] + amount[len(v):]
	params := "[{\"FROM\":\"" + from + "\", \"TO\":\"" + contract + "\", \"VALUE\":\"0x0\", \"DATA\":\"" + data + "\"},\"" + password + "\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_sendTxWithPassword\", \"PARAMS\":" + params + ", \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmSendTxWithPasswordTokenResp
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
func Rpc_signTxWithPassword(from string, to string, value float64, password string) (DDM_SignTxInfo, error) {
	var rtn DDM_SignTxInfo

	if len(from) != 42 {
		err := errors.New("{\"from: non-regular address.\"}")
		return rtn, err
	}
	if len(to) != 42 {
		err := errors.New("{\"to: non-regular address.\"}")
		return rtn, err
	}

	val := util.ToWei(value, 18)
	v := util.BigInt2HexString(val)

	gas, gasPrice, nonce, err := getSignTxInfo(from)
	if err != nil {
		return rtn, err
	}

	params := "[{\"FROM\":\"" + from + "\", \"TO\":\"" + to + "\", \"VALUE\":\"0x" + v + "\", \"GAS\":\"" + gas + "\", \"GASPRICE\":\"" + gasPrice + "\", \"NONCE\":\"" + nonce + "\"},\"" + password + "\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_signTxWithPassword\", \"PARAMS\":" + params + ", \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmSignTxWithPasswordResp
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
func Rpc_sendRawTx(raw string) (string, error) {
	var rtn string

	params := "[\"" + raw + "\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_sendRawTx\", \"PARAMS\":" + params + ", \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmSendRawTxResp
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
func Rpc_gasPrice() (*big.Int, error) {
	rtn := big.NewInt(0)

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_gasPrice\", \"PARAMS\":[], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetPriceResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		val := bResp.Result[2:]
		rtn, _ = rtn.SetString(val, 16) // hex string to big.Int
	}

	return rtn, nil
}
func Rpc_getContractToken(contract string) (string, error) {
	var rtn string = ""

	params := "[{\"DATA\":\"0x95d89b41\", \"TO\":\"" + contract + "\"},\"latest\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_call\", \"PARAMS\":" + params + ", \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetContractTokenResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		if len(bResp.Result) > 2 {
			rtn = util.Hex2AsciiString(bResp.Result[130:])
			if rtn == "" {
				return rtn, errors.New("{\"failure to got the symbol.\"}")
			}
			rtn = strings.Replace(rtn, "\u0000", "", -1)
		}
	}

	return rtn, nil
}
func Rpc_getContractName(contract string) (string, error) {
	var rtn string = ""

	params := "[{\"DATA\":\"0x06fdde03\", \"TO\":\"" + contract + "\"},\"latest\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_call\", \"PARAMS\":" + params + ", \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmGetContractNameResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		if len(bResp.Result) > 2 {
			rtn = util.Hex2AsciiString(bResp.Result[130:])
			if rtn == "" {
				return rtn, errors.New("{\"failure to got the symbol.\"}")
			}
			rtn = strings.Replace(rtn, "\u0000", "", -1)
		}
	}

	return rtn, nil
}
func Rpc_blockNumber() (int64, error) {
	var rtn int64 = 0

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"ddm_blockNumber\", \"PARAMS\":[], \"ID\":101}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ddmBlockNumberResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		val := bResp.Result[2:]
		rtn, err = strconv.ParseInt(val, 16, 64)
		if err != nil {
			return rtn, err
		}
	}

	return rtn, nil
}
