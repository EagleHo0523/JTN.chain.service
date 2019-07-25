package mqtt

import (
	"fmt"
	"time"

	util "../util"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func CreateConnect(mqSetting util.MqttSetting) (mqtt.Client, error) {
	var client mqtt.Client
	opts := mqtt.NewClientOptions()

	if mqSetting.Env == "main" {
		// fmt.Println("#1, ", mqSetting.Url)
		opts.AddBroker(mqSetting.Url)
		opts.SetUsername(mqSetting.Username)
		opts.SetPassword(mqSetting.Password)
	} else {
		opts.AddBroker(fmt.Sprintf("%s://%s:%s", "tcp", "localhost", "1883"))
	}

	// opts.SetClientID("test-clientID") // Multiple connections should use different clientID for each connection, or just leave it blank
	opts.SetKeepAlive(time.Second * time.Duration(60))
	// If lost connection, reconnect again
	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		util.FailOnError(err, "Connection lost :")
	})

	// connect to broker
	client = mqtt.NewClient(opts)

	return client, nil
}

func Send(client mqtt.Client, topic string, data interface{}) {
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		util.FailOnError(token.Error(), "Fail to connect broker,")
	}

	body := data.([]byte)

	// publish to topic
	token = client.Publish(topic, byte(2), false, body) // QoS 2 較占用頻寬, 但保證送達, 且只送一次
	if token.Wait() && token.Error() != nil {
		util.FailOnError(token.Error(), "Fail to publish,")
	}
	// fmt.Println("#2")
}
