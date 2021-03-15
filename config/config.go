package config

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/vivian-tangle/vivian-node/tools"
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

	configDir = "config.json"
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
	viper.AddConfigPath(".")
	viper.SetConfigFile(configDir)
	// Load default configurations
	viper.SetDefault("network", DefaultNetwork)
	viper.SetDefault("databasePath", DefaultDatabasePath)
	viper.SetDefault("node", DefaultNode)
	viper.SetDefault("zmqSocket", DefaultZmqSocket)

	// Load configurations from the configuration file
	err := viper.ReadInConfig()
	tools.HandleErr(err)

	// Unmarshal config
	err = viper.Unmarshal(&c)
	tools.HandleErr(err)

	fmt.Println("Configuration loaded")
}
