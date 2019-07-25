package hdwallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
)

type HDWallet struct {
	Params *chaincfg.Params

	// master key options
	Mnemonic string
	Password string

	RootPrivateKey string
	RootPublicKey  string
	ETH            AccountInfo
	BTC            AccountInfo
	USDT           AccountInfo
}
type AccountInfo struct {
	Key *hdkeychain.ExtendedKey
	// RootPrivateKey string
	// RootPublicKey  string
	Account    string
	PrivateKey string
	PublicKey  string
}

type dataBase struct {
	Url      string `json:"Url"`
	DBname   string `json:"DBname"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type settingValue struct {
	Env string `json:"Env"`
	Url string `json:"Url"`
	Key string `json:"Key"`
	db  dataBase
}
type envSetting struct {
	DDM settingValue
	ETH settingValue
	BTC settingValue
}

// NewMnemonic ...
func NewMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}
func NewWalletAccount(mnemonic string, passphrase string) (HDWallet, error) {
	var rtn HDWallet

	env, err := getSystemParameters("BTC")
	if err != nil {
		return rtn, err
	}
	switch env.Env {
	case "main":
		rtn.Params = &chaincfg.MainNetParams
	case "test":
		rtn.Params = &chaincfg.TestNet3Params
	case "reg":
		rtn.Params = &chaincfg.RegressionNetParams
	}

	if mnemonic == "" {
		return rtn, errors.New("{\"mnemonic is required\"}")
	} else {
		rtn.Mnemonic = mnemonic
	}
	if passphrase == "" {
		return rtn, errors.New("{\"passhrase is required\"}")
	} else {
		rtn.Password = passphrase
	}

	seed := bip39.NewSeed(rtn.Mnemonic, rtn.Password)
	privKey, err := hdkeychain.NewMaster(seed, rtn.Params) // *ExtendedKey, error // MainNetParams, TestNet3Params, RegressionNetParams, SimNetParams
	if err != nil {
		return rtn, err
	}
	rootPrivKey := fmt.Sprintf("%v", privKey)
	pubKey, _ := privKey.Neuter() // for uuid
	rootPubKey := fmt.Sprintf("%v", pubKey)

	ethAccount, err := newEthAccount(rootPrivKey)
	btcAccount, err := newBtcAccount(rootPrivKey, rtn.Params)
	usdtAccount, err := newUsdtAccount(rootPrivKey, rtn.Params)

	rtn.RootPrivateKey = rootPrivKey
	rtn.RootPublicKey = rootPubKey
	rtn.ETH = ethAccount
	rtn.BTC = btcAccount
	rtn.USDT = usdtAccount

	return rtn, nil
}
func LoopEthAccount(mnemonic string, passphrase string) {
	seed := bip39.NewSeed(mnemonic, passphrase)
	privKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams) // *ExtendedKey, error // MainNetParams, TestNet3Params, RegressionNetParams, SimNetParams
	if err != nil {
		return
	}
	rootPrivKey := fmt.Sprintf("%v", privKey)
	// pubKey, _ := privKey.Neuter() // for uuid
	// rootPubKey := fmt.Sprintf("%v", pubKey)

	for i := 31550941; i <= 4294967295; i++ {
		ethKey, err := newKey(rootPrivKey, 44, 60, 0, 0, uint32(i))
		if err != nil {
			return
		}
		privateKey, err := ethKey.ECPrivKey()
		if err != nil {
			return
		}
		privateKeyECDSA := privateKey.ToECDSA()
		// privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)

		publicKey := privateKeyECDSA.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return
		}
		// publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

		addr := crypto.PubkeyToAddress(*publicKeyECDSA)
		addrHex := addr.Hex()
		if strings.Index(addrHex, "0xEe9b") > -1 {
			fmt.Println(i, " account: ", addrHex)
		}
	}
}

func encryptKey(text string, passphrase string) (string, error) {
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

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	seal := gcm.Seal(nonce, nonce, data, nil)

	return base58.Encode(seal), nil
}
func decryptKey(text string, passphrase string) (string, error) {
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
func newKey(master string, purpose uint32, coinType uint32, account uint32, change uint32, addressIndex uint32) (*hdkeychain.ExtendedKey, error) {
	var rtn *hdkeychain.ExtendedKey

	masterKey, err := hdkeychain.NewKeyFromString(master)
	if err != nil {
		return rtn, err
	}

	//   m/44'
	acct44, err := masterKey.Child(hdkeychain.HardenedKeyStart + purpose)
	if err != nil {
		return rtn, err
	}
	acct44Ext, err := acct44.Child(hdkeychain.HardenedKeyStart + coinType)
	if err != nil {
		return rtn, err
	}
	acct44ExtExt, err := acct44Ext.Child(hdkeychain.HardenedKeyStart + account)
	if err != nil {
		return rtn, err
	}
	acct44ExtExtExt, err := acct44ExtExt.Child(change)
	if err != nil {
		return rtn, err
	}
	acct44ExtExtExtExt, err := acct44ExtExtExt.Child(addressIndex)
	if err != nil {
		return rtn, err
	}
	rtn = acct44ExtExtExtExt

	return rtn, nil
}
func newBtcAccount(master string, params *chaincfg.Params) (AccountInfo, error) {
	var rtn AccountInfo

	btcKey, err := newKey(master, 44, 0, 0, 0, 0)
	if err != nil {
		return rtn, err
	}

	addr, _ := btcKey.Address(params)

	privateKey, err := btcKey.ECPrivKey()
	if err != nil {
		return rtn, err
	}

	publicKey, err := btcKey.ECPubKey()
	if err != nil {
		return rtn, err
	}
	publicKeyBytes := publicKey.SerializeCompressed()

	wif, _ := btcutil.NewWIF(privateKey, params, true)

	rtn.Key = btcKey
	rtn.PrivateKey = wif.String()
	rtn.PublicKey = hex.EncodeToString(publicKeyBytes)
	rtn.Account = addr.String()

	return rtn, nil
}
func newEthAccount(master string) (AccountInfo, error) {
	var rtn AccountInfo

	ethKey, err := newKey(master, 44, 60, 0, 0, 0)
	if err != nil {
		return rtn, err
	}

	privateKey, err := ethKey.ECPrivKey()
	if err != nil {
		return rtn, err
	}
	privateKeyECDSA := privateKey.ToECDSA()
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return rtn, errors.New("{\"failed ot get public key\"}")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	addr := crypto.PubkeyToAddress(*publicKeyECDSA)

	rtn.Key = ethKey
	rtn.PrivateKey = hexutil.Encode(privateKeyBytes)[2:]
	rtn.PublicKey = hexutil.Encode(publicKeyBytes)[2:]
	rtn.Account = addr.Hex()

	return rtn, nil
}
func newUsdtAccount(master string, params *chaincfg.Params) (AccountInfo, error) {
	var rtn AccountInfo

	usdtKey, err := newKey(master, 44, 0, 0, 0, 0) // same with btc
	if err != nil {
		return rtn, err
	}

	addr, _ := usdtKey.Address(params)

	privateKey, err := usdtKey.ECPrivKey()
	if err != nil {
		return rtn, err
	}

	publicKey, err := usdtKey.ECPubKey()
	if err != nil {
		return rtn, err
	}
	publicKeyBytes := publicKey.SerializeCompressed()

	wif, _ := btcutil.NewWIF(privateKey, params, true)

	rtn.Key = usdtKey
	rtn.PrivateKey = wif.String()
	rtn.PublicKey = hex.EncodeToString(publicKeyBytes)
	rtn.Account = addr.String()

	return rtn, nil
}

func getSystemParameters(method string) (settingValue, error) {
	var rtn settingValue
	jsonFile, err := os.Open("env.json")
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
	case "DDM":
		rtn = envs.DDM
	case "ETH":
		rtn = envs.ETH
	case "BTC":
		rtn = envs.BTC
	}
	return rtn, nil
}
