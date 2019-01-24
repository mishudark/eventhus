package mosquitto

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/mishudark/eventhus/v2"
	"log"
	"os"
)

var info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

var fatal = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)

// Client nats
type Client struct {
	options *MQTT.ClientOptions
	client  MQTT.Client
}

//MqttDefaultPort is the default port
const MqttDefaultPort = 1883

//MqttDefaultHost is the default host
const MqttDefaultHost = "localhost"

//MqttDefaultMethod is the default method
const MqttDefaultMethod = "tcp"

//MqttDefaultClientId is the default method
const MqttDefaultClientId = "cqrs-es"

//NewClient create a new client with default parameters
func NewClient() (*Client, error) {
	return NewClientWithPort(MqttDefaultMethod, MqttDefaultHost, MqttDefaultPort, MqttDefaultClientId)
}

// defaultPublishHandler is called when there is a message that is matching no other subscriber
var defaultPublishHandler MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

//NewClientWithPort create a new client with options
func NewClientWithPort(method string, host string, port int, clientID string) (*Client, error) {
	var d Client

	d.options = MQTT.NewClientOptions()
	brokerURL := fmt.Sprintf("%s://%s:%d", method, host, port)
	d.options.AddBroker(brokerURL)
	d.options.SetClientID(clientID)
	d.options.SetDefaultPublishHandler(defaultPublishHandler)

	info.Printf("Created Client for broker %s", brokerURL)

	return &d, nil
}

// Publish a event
func (c *Client) Publish(event eventhus.Event, bucket, subset string) error {

	info.Println("Publish Begin")

	c.client = MQTT.NewClient(c.options)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	defer c.client.Disconnect(5000)

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	subj := bucket + "/" + subset
	token := c.client.Publish(subj, 0, false, msg)
	token.Wait()

	info.Println("Publish End")

	return token.Error()
}
