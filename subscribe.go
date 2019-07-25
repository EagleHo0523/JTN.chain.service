package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	util "./util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type responseBlockNumber struct {
	BlockNumber string
}

func main() {
	go Receive("blockchain/ddmx/listenblocknumber") // blockchain/+/listenblocknumber
	Receive("blockchain/ddmx/listentx/Host")
}

func Receive(topic string) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://message.shopeebuy.com:1883") //opts.AddBroker(fmt.Sprintf("%s://%s:%s", "tcp", "localhost", "1883")) // message.shopeebuy.com
	opts.SetUsername("jutaibc")
	opts.SetPassword("1qaz2wsx@BC")
	// opts.SetClientID("test-clientID") // Multiple connections should use different clientID for each connection, or just leave it blank
	opts.SetKeepAlive(time.Second * time.Duration(60))
	// If lost connection, reconnect again
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		util.FailOnError(err, "Connection lost :")
		if client.IsConnected() {
			client.Disconnect(500)
		}
		client = mqtt.NewClient(opts)
		connect(client)
		subscribe(client, topic)
		fmt.Println("#2")
	})

	client := mqtt.NewClient(opts)

	connect(client)
	subscribe(client, topic)
	fmt.Println("#3")

	fmt.Println("Start to subcribe...")

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("unknown panic error, try to recover connection,", r)
			connect(client)
			subscribe(client, topic)
			fmt.Println("#4")
		}
	}()

	select {}
}

func connect(client mqtt.Client) {
	for {
		token := client.Connect()
		if token.Wait() && token.Error() != nil {
			util.FailOnError(token.Error(), "Fail to connect broker,")
			time.Sleep(5 * time.Second)

			fmt.Println("Retry the connection...")
			continue
		} else {
			fmt.Println("Reconnection successful!")
			break
		}
	}
}

func subscribe(client mqtt.Client, topic string) {
	for {
		token := client.Subscribe(topic, byte(2), onIncomingDataReceived)
		if token.Wait() && token.Error() != nil {
			util.FailOnError(token.Error(), "Fail to sub...")
			time.Sleep(5 * time.Second)

			fmt.Println("Retry to subscribe...")
			continue
		} else {
			fmt.Println("Subscribe successful!")
			break
		}
	}
}

func onIncomingDataReceived(client mqtt.Client, message mqtt.Message) {
	topic := message.Topic()
	if strings.Index(topic, "listentx") >= 0 {
		var backTx util.BackTxMessage
		err := json.Unmarshal(message.Payload(), &backTx)
		util.FailOnError(err, "Failed to Unmarshal JSON")
		fmt.Println("Topic: ", message.Topic())
		fmt.Println("BlockNumber: ", backTx.BlockNumber)
		fmt.Println("TxHash: ", backTx.TxHash)
		fmt.Println("From: ", backTx.From)
		fmt.Println("To: ", backTx.To)
		fmt.Println("Contract: ", backTx.Contract)
		fmt.Println("Amount: ", backTx.Amount)
		fmt.Println("Symbol: ", backTx.Symbol)
	}
	if strings.Index(topic, "listenblocknumber") >= 0 {
		var blockNum responseBlockNumber
		err := json.Unmarshal(message.Payload(), &blockNum)
		util.FailOnError(err, "Failed to Unmarshal JSON")
		fmt.Println("Datetime: ", time.Now())
		fmt.Println("BlockNumber: ", blockNum.BlockNumber)
	}
}

func UnSubscribe(client mqtt.Client) {
	//unsubscribe from /go-mqtt/sample
	if token := client.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		util.FailOnError(token.Error(), "Fail to unsubscribe,")
		os.Exit(1)
	}

	client.Disconnect(250)
}
