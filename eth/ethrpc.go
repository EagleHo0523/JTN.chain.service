package eth

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	util "../util"
)

type ETH_BlockInfo struct {
	Difficulty          string
	ExtraData           string
	GasLimit            string
	GasUsed             string
	Hash                string
	LogsBloom           string
	Miner               string
	MixHash             string
	Nonce               string
	Number              string
	ParentHash          string
	ReceiptsRoot        string
	Sha3Uncles          string
	Size                string
	StateRoot           string
	Timestamp           string
	TotalDifficulty     string
	TransactionReceipts []ETH_TransactionReceipts `json:",omitempty"`
	Transactions        []ETH_TransactionInfo     `json:",omitempty"`
	TransactionsRoot    string
	// UNCLES
}
type ETH_TransactionReceipts struct {
	CONTRACTADDRESS string `json:",omitempty"`
	GASUSED         string
	LOGS            []ETH_ReceiptLogs `json:",omitempty"`
	STATUS          string
	TRANSACTIONHASH string
}
type ETH_ReceiptLogs struct {
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
type ETH_TransactionInfo struct {
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
type ETH_TransactionReceiptInfo struct {
	BLOCKHASH         string
	BLOCKNUMBER       string
	CONTRACTADDRESS   string `json:",omitempty"`
	CUMCLATIVEGASUSED string
	FROM              string
	GASUSED           string
	LOGS              []string `json:",omitempty"`
	LOGSBLOOM         string
	STATUS            string
	TO                string
	TRANSACTIONHASH   string
	TRANSACTIONINDEX  string
}

type ethError struct {
	Code    int64
	Message string
}
type ethGetPriceResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethAccountsResp struct {
	Jsonrpc string
	Id      int
	Result  []string `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethNewAccountResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethGetBlockByNumberResp struct {
	Jsonrpc string
	Id      int
	Result  ETH_BlockInfo `json:",omitempty"`
	Error   ethError      `json:",omitempty"`
}
type ethGetBlockByHashResp struct {
	Jsonrpc string
	Id      int
	Result  ETH_BlockInfo `json:",omitempty"`
	Error   ethError      `json:",omitempty"`
}
type ethGetBlockTransactionCountByHashResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethGetBlockTransactionCountByNumberResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethGetTransactionByHashResp struct {
	Jsonrpc string
	Id      int
	Result  ETH_TransactionInfo `json:",omitempty"`
	Error   ethError            `json:",omitempty"`
}
type ethGetTransactionReceiptResp struct {
	Jsonrpc string
	Id      int
	Result  ETH_TransactionReceiptInfo `json:",omitempty"`
	Error   ethError                   `json:",omitempty"`
}
type ethGetBalanceResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethSendTransactionResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethSendTransactionTokenResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethGetContractTokenResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethGetContractNameResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}
type ethBlockNumberResp struct {
	Jsonrpc string
	Id      int
	Result  string   `json:",omitempty"`
	Error   ethError `json:",omitempty"`
}

func processError(m interface{}) error {
	e := m.(ethError)
	code := fmt.Sprintf("%v", e.Code)
	message := e.Message
	err := errors.New("{\"CODE\":\"" + code + "\",\"MESSAGE\":\"" + message + "\"}")
	return err
}

func Rpc_gasPrice() (*big.Int, error) {
	rtn := big.NewInt(0)

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_gasPrice\", \"PARAMS\":[], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetPriceResp
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
func Rpc_getAccounts() ([]string, error) {
	var rtn []string

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_accounts\", \"PARAMS\":[], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethAccountsResp
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
func Rpc_newAccount(passpharse string) (string, error) {
	var rtn string = ""

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"personal_newAccount\", \"PARAMS\":[\"" + passpharse + "\"], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethNewAccountResp
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
func Rpc_getBlockByNumber(num string) (ETH_BlockInfo, error) {
	var rtn ETH_BlockInfo

	// 分析 num 字串, 是否為數字; 如為文字, 是否為特定字串

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_getBlockByNumber\", \"PARAMS\":[\"" + num + "\",true], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetBlockByNumberResp
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
func Rpc_getBlockByHash(hashValue string) (ETH_BlockInfo, error) {
	var rtn ETH_BlockInfo

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_getBlockByHash\", \"PARAMS\":[\"" + hashValue + "\",true], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetBlockByHashResp
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
func Rpc_getBlockTransactionCountByHash(hashValue string) (int64, error) {
	var rtn int64

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_getBlockTransactionCountByHash\", \"PARAMS\":[\"" + hashValue + "\"], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetBlockTransactionCountByHashResp
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
func Rpc_getBlockTransactionCountByNumber(num int64) (int64, error) {
	var rtn int64

	hexStr := strconv.FormatInt(num, 16) //int64 to hex string

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_getBlockTransactionCountByNumber\", \"PARAMS\":[\"0x" + hexStr + "\"], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetBlockTransactionCountByNumberResp
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
func Rpc_getTransactionByHash(hashValue string) (ETH_TransactionInfo, error) {
	var rtn ETH_TransactionInfo

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_getTransactionByHash\", \"PARAMS\":[\"" + hashValue + "\"], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetTransactionByHashResp
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
func Rpc_getTransactionReceipt(hashValue string) (ETH_TransactionReceiptInfo, error) {
	var rtn ETH_TransactionReceiptInfo

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_getTransactionReceipt\", \"PARAMS\":[\"" + hashValue + "\"], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetTransactionReceiptResp
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
func Rpc_getBalance(account string) (float64, error) {
	var rtn float64

	params := "[\"" + account + "\",\"latest\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_getBalance\", \"PARAMS\":" + params + ", \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetBlockTransactionCountByHashResp
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
func Rpc_getBalanceToken(account string, contract string) (float64, error) {
	var rtn float64

	params := "[{\"DATA\":\"0x70a08231000000000000000000000000" + account[2:] + "\", \"TO\":\"" + contract + "\"},\"latest\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_call\", \"PARAMS\":" + params + ", \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetBalanceResp
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
func Rpc_sendTransaction(from string, to string, value float64, password string) (string, error) {
	var rtn string = ""

	if len(from) != 42 {
		err := errors.New("{\"from: non-regular address.\"}")
		return "", err
	}
	if len(to) != 42 {
		err := errors.New("{\"to: non-regular address\"}")
		return "", err
	}

	val := util.ToWei(value, 18)
	v := util.BigInt2HexString(val)

	params := "[{\"FROM\":\"" + from + "\", \"TO\":\"" + to + "\", \"VALUE\":\"0x" + v + "\"},\"" + password + "\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"personal_sendTransaction\", \"PARAMS\":" + params + ", \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethSendTransactionResp
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
func Rpc_sendTransactionToken(from string, to string, value float64, password string, contract string) (string, error) {
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
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"personal_sendTransaction\", \"PARAMS\":" + params + ", \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethSendTransactionTokenResp
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
func Rpc_getContractToken(contract string) (string, error) {
	var rtn string = ""

	params := "[{\"DATA\":\"0x95d89b41\", \"TO\":\"" + contract + "\"},\"latest\"]"
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_call\", \"PARAMS\":" + params + ", \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethSendTransactionTokenResp
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
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_call\", \"PARAMS\":" + params + ", \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetContractNameResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		if len(bResp.Result) > 2 {
			fmt.Println("bResp.Result: ", bResp.Result)
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

	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_blockNumber\", \"PARAMS\":[], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethBlockNumberResp
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

func abc() (int64, error) {
	var rtn int64 = 0
	args := "{\"JSONRPC\":\"2.0\", \"METHOD\":\"eth_getBlockByNumber\", \"PARAMS\":[\"latest\",true], \"ID\":33571}"

	resp, err := util.RpcPost(url, args)
	if err != nil {
		return rtn, err
	}

	var bResp ethGetBlockByNumberResp
	err = json.Unmarshal([]byte(resp), &bResp)
	if err != nil {
		return rtn, err
	}

	if bResp.Error.Code != 0 {
		return rtn, processError(bResp.Error)
	} else {
		val := bResp.Result.Number[2:]
		rtn, err = strconv.ParseInt(val, 16, 64)
		if err != nil {
			return rtn, err
		}
	}
	return rtn, nil
}
