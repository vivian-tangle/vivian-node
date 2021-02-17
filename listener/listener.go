package listener

import (
	"fmt"
	"strings"

	zmq "github.com/pebbe/zmq4"
	"github.com/vivian-tangle/vivian-node/config"
)

// Listener is the struct for storing the information of a ZMQ listener
type Listener struct {
	Config *config.Config
}

// Listen starts listening to ZMQ and handle the transactions
func (listener *Listener) Listen() {
	client, _ := zmq.NewSocket(zmq.SUB)

	// Make sure the connection is closed after stopping the program
	defer client.Close()

	client.Connect(listener.Config.ZmqSocket)

	// Subscribe to both tx and sn (confirmed tx) topics
	client.SetSubscribe("tx")
	client.SetSubscribe("sn")

	// Keep looping for messages
	for {
		msg, _ := client.RecvMessage(0)
		for _, str := range msg {
			// Split the fields by the space character
			words := strings.Fields(str)

			if words[0] == "tx" {
				fmt.Println("New transaction: ", words[1])
			}
			if words[0] == "sn" {
				fmt.Println("Confirmed transaction: ", words[2], "for milestone", words[1])
			}
		}
	}
}
