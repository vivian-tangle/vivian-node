package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	// DefaultNetwork is the default parameter for choosing IOTA network (mainnet/devnet)
	DefaultNetwork = "devnet"
	// DefaultDatabasePath is directory for storing database files
	DefaultDatabasePath = "db"
	// DefaultNode is the default node for connecting to IOTA network
	DefaultNode = "https://nodes.devnet.iota.org"
	// DefaultZmqSocket is the default connection address and port of ZMQ
	DefaultZmqSocket = "tcp://zmq.devnet.iota.org:5556"
	configDir        = "config.json"
)

// Config is the struct for storing config parameters
type Config struct {
	Network      string
	DatabasePath string
	Node         string
	ZmqSocket    string
}

// LoadConfig loads the configures from default config json
func (c *Config) LoadConfig() {
	// Load default configurations
	c.Network = DefaultNetwork
	c.DatabasePath = DefaultDatabasePath
	c.Node = DefaultNode
	c.ZmqSocket = DefaultZmqSocket

	configJSON, err := os.Open(configDir)
	if err != nil {
		fmt.Printf("Cannot load %s, using default configure values...", configDir)
		return
	}
	defer configJSON.Close()

	// Read the opened config json as a byte array
	byteValue, _ := ioutil.ReadAll(configJSON)
	var data map[string]string
	err = json.Unmarshal(byteValue, &data)
	handleErr(err)
	if val, ok := data["network"]; ok {
		c.Network = val
	}
	if val, ok := data["databasePath"]; ok {
		c.DatabasePath = val
	}
	if val, ok := data["node"]; ok {
		c.Node = val
	}
	if val, ok := data["zmqSocket"]; ok {
		c.Node = val
	}

	fmt.Println("Configuration loaded")
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
