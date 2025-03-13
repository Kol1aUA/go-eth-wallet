package wallet

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

const walletsFile = "wallets.json"

var (
	mu             sync.Mutex
	defaultWallet  string
	defaultKeyFile = "default_wallet.txt" // Stores the default wallet name
)

// WalletData represents the structure of stored wallets
type WalletData struct {
	Wallets map[string]string `json:"wallets"` // Map of wallet name → private key
}

// LoadWallets loads wallets from file
func LoadWallets() (*WalletData, error) {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.ReadFile(walletsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &WalletData{Wallets: make(map[string]string)}, nil // Return empty wallets if file doesn't exist
		}
		return nil, err
	}

	var data WalletData
	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// SaveWallets saves wallets to file
func SaveWallets(wallets *WalletData) error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(wallets, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(walletsFile, data, 0600)
}

// SetDefaultWallet saves the default wallet name
func SetDefaultWallet(walletName string) error {
	mu.Lock()
	defer mu.Unlock()

	return os.WriteFile(defaultKeyFile, []byte(walletName), 0600)
}

// GetDefaultWallet retrieves the default wallet name
func GetDefaultWallet() (string, error) {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(defaultKeyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("no default wallet set")
		}
		return "", err
	}

	return string(data), nil
}
