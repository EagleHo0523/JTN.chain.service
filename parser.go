package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type dataBase struct {
	Url      string `json:"Url"`
	DBname   string `json:"DBname"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}
type settingValue struct {
	Url string `json:"Url"`
	Key string `json:"Key"`
	db  dataBase
}
type envSetting struct {
	DDM settingValue
	ETH settingValue
	BTC settingValue
}

func parser() {
	env, err := getSystemParameters("DDM")
	if err != nil {
		fmt.Println("{\"" + err.Error() + "\"}")
		return
	} else {
		fmt.Println(env.Url)
		fmt.Println(env.Key)
	}

}

func getSystemParameters(method string) (settingValue, error) {
	var rtn settingValue
	jsonFile, err := os.Open("/home/ubuntu/chain_service_go/env.json")
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
