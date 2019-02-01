package main

import "github.com/spf13/viper"

const (
	cfgNodeHost     = "Node.Host"
	cfgNodeUser     = "Node.User"
	cfgNodePassword = "Node.Password"
)

// Miner represents the Miner section of the configuration
type Node struct {
}

func newNode() *Node {
	miner := new(Node)

	viper.SetDefault(cfgNodeHost, "localhost:18331")
	viper.SetDefault(cfgNodeUser, "equibit")
	viper.SetDefault(cfgNodePassword, "equibit")

	return miner
}
