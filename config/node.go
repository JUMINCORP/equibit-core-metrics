package config

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

func (*Node) Host() string {
	return viper.Get(cfgNodeHost).(string)
}

func (*Node) User() string {
	return viper.Get(cfgNodeUser).(string)
}

func (*Node) Password() string {
	return viper.Get(cfgNodePassword).(string)
}
