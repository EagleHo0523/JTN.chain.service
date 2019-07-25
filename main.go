package main

import (
	// w "./hdwallet"
	"log"
	"net/http"

	ddmx "./ddmx"
	// eth "./eth"
	r "./router"
)

func main() {
	go ddmx.ProcessTxBack()
	// go eth.ProcessTxBack()
	// go btc.ProcessTxBack()
	apiService()

	// mnemonic, err := w.NewMnemonic()
	// if err != nil {
	// 	return
	// }
	// mnemonic := "orbit crumble stand output swift solar orange assist rescue share cherry allow"
	// mnemonic := "antenna genuine only private urge apart amount lawn educate amazing wreck market"
	// mnemonic := "toilet venture gain captain coffee dinner eyebrow wagon draft figure lottery absorb"
	// passphrase := "Eagle@0523"
	// w.LoopEthAccount(mnemonic, passphrase)
	// wallet, err := w.NewWalletAccount(mnemonic, passphrase)
	// if err != nil {
	// 	return
	// }
	// fmt.Println("mnemonic: ", mnemonic)
	// fmt.Println("RootPrivateKey: ", wallet.RootPrivateKey)
	// fmt.Println("RootPublicKey: ", wallet.RootPublicKey)
	// fmt.Println("===================================================================================================================")
	// fmt.Println("BTC private key: ", wallet.BTC.PrivateKey)
	// fmt.Println("BTC public key: ", wallet.BTC.PublicKey)
	// fmt.Println("BTC account: ", wallet.BTC.Account)
	// fmt.Println("===================================================================================================================")
	// fmt.Println("ETH private key: ", wallet.ETH.PrivateKey)
	// fmt.Println("ETH public key: ", wallet.ETH.PublicKey)
	// fmt.Println("ETH account: ", wallet.ETH.Account)
	// fmt.Println("===================================================================================================================")
	// fmt.Println("USDT private key: ", wallet.USDT.PrivateKey)
	// fmt.Println("USDT public key: ", wallet.USDT.PublicKey)
	// fmt.Println("USDT account: ", wallet.USDT.Account)
	// fmt.Println("===================================================================================================================")
	// e, err := w.EncryptKey(wallet.RootPrivateKey, passphrase)
	// if err != nil {
	// 	fmt.Println("err #1: ", err)
	// }
	// fmt.Println("Encrypt string: ", e)
	// fmt.Println("Encrypt len: ", len(e))
	// fmt.Println("===================================================================================================================")
	// d, err := w.DecryptKey(e, passphrase)
	// if err != nil {
	// 	fmt.Println("err #2: ", err)
	// }
	// fmt.Println("Decrypt string: ", d)
	// fmt.Println("same before: ", strings.Compare(wallet.RootPrivateKey, d))
}

func apiService() {
	router := r.NewRouter()
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Chain service listen and serve: ", err)
	}
}
