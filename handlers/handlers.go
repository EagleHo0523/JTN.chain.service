package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	btc "../btc"
	ddm "../ddmx"
	eth "../eth"
)

func DDMX(w http.ResponseWriter, r *http.Request) {
	method, args, err := processRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		resp, err := ddm.ProcessMethod(method, args)
		if err != nil {
			http.Error(w, "{\""+err.Error()+"\"}", http.StatusBadRequest)
			return
		} else {
			// w.Header().Set("Content-Type", "application/json")
			processResponse(resp, w)
		}
	}
}
func ETH(w http.ResponseWriter, r *http.Request) {
	method, args, err := processRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		resp, err := eth.ProcessMethod(method, args)
		if err != nil {
			http.Error(w, "{\""+err.Error()+"\"}", http.StatusBadRequest)
			return
		} else {
			// w.Header().Set("Content-Type", "application/json")
			processResponse(resp, w)
		}
	}
}
func BTC(w http.ResponseWriter, r *http.Request) {
	method, args, err := processRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		resp, err := btc.ProcessMethod(method, args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			// w.Header().Set("Content-Type", "application/json")
			processResponse(resp, w)
		}
	}
}

func processRequest(r *http.Request) (string, interface{}, error) {
	var method string = ""
	var params interface{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return method, params, err
	}

	var f interface{}
	decoder := json.NewDecoder(bytes.NewBuffer(body))
	decoder.UseNumber() // 此处能够保证bigint的精度
	decoder.Decode(&f)

	m := f.(map[string]interface{})
	method = fmt.Sprintf("%v", m["METHOD"])
	tmp := m["PARAMS"]
	params = tmp.(interface{})

	return method, params, nil
}
func processResponse(resp interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
	}
}
