package ethclient

import (
	"errors"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	clientInstance *ethclient.Client
	once           sync.Once
)

// GetClient returns a singleton Ethereum client
func GetClient() (*ethclient.Client, error) {
	var err error
	once.Do(func() {
		rpcURL := os.Getenv("ETH_RPC_URL")
		if rpcURL == "" {
			err = errors.New("Ethereum RPC URL not set. Please set the ETH_RPC_URL environment variable")
			return
		}

		clientInstance, err = ethclient.Dial(rpcURL)
	})

	return clientInstance, err
}
