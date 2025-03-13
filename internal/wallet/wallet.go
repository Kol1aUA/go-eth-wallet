package wallet

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"go-eth-wallet/internal/ethclient"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var walletPrivateKey = "WALLET_PRIVATE_KEY"

type Transaction struct {
	Hash      string `json:"hash"`
	From      string `json:"from"`
	To        string `json:"to"`
	Value     string `json:"value"`
	Timestamp string `json:"timeStamp"`
}

// CreateNewWallet creates a new Ethereum wallet and stores it under a given name
func CreateNewWallet(walletName string) (string, error) {
	// Generate a new Ethereum key pair
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", err
	}

	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))
	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	// Load existing wallets
	wallets, err := LoadWallets()
	if err != nil {
		return "", err
	}

	// Store wallet
	wallets.Wallets[walletName] = privateKeyHex
	err = SaveWallets(wallets)
	if err != nil {
		return "", errors.New("failed to save wallet")
	}

	// Set as default wallet if it's the first wallet
	if len(wallets.Wallets) == 1 {
		SetDefaultWallet(walletName)
	}

	return address, nil
}

// ImportWallet allows importing an existing private key with a given name
func ImportWallet(walletName, privateKeyHex string) error {
	_, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return errors.New("invalid private key format")
	}

	// Load wallets
	wallets, err := LoadWallets()
	if err != nil {
		return err
	}

	// Store wallet
	wallets.Wallets[walletName] = privateKeyHex
	err = SaveWallets(wallets)
	if err != nil {
		return errors.New("failed to save imported wallet")
	}

	return nil
}

// GetAddress retrieves the Ethereum address of the default wallet
func GetAddress() (string, error) {
	wallets, err := LoadWallets()
	if err != nil {
		return "", err
	}

	defaultWallet, err := GetDefaultWallet()
	if err != nil {
		return "", err
	}

	privateKeyHex, exists := wallets.Wallets[defaultWallet]
	if !exists {
		return "", errors.New("default wallet not found")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", errors.New("invalid private key format")
	}

	return crypto.PubkeyToAddress(privateKey.PublicKey).Hex(), nil
}

func GetBalance() (*big.Float, error) {
	address, err := GetAddress()
	if err != nil {
		return nil, err
	}

	client, err := ethclient.GetClient()
	if err != nil {
		return nil, err
	}

	account := common.HexToAddress(address)

	balance, err := client.BalanceAt(context.Background(), account, nil)

	if err != nil {
		return nil, errors.New("failed to retrieve balance")
	}

	ethBalance := new(big.Float).SetInt(balance)
	ethBalance.Quo(ethBalance, big.NewFloat(1e18))

	return ethBalance, nil
}

func SendETH(recipient string, amount string) (string, error) {

	// Ensure a wallet is available
	err := EnsureWalletExists()
	if err != nil {
		return "", err
	}

	// Retrieve private key
	privateKeyHex := GetPrivateKeyFromEnv()
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return "", errors.New("invalid private key format")
	}
	// Convert recipient to Ethereum address
	toAddress := common.HexToAddress(recipient)

	// Connect to Ethereum node
	client, err := ethclient.GetClient()
	if err != nil {
		return "", errors.New("failed to connect to Ethereum node")
	}

	// Get sender's address
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)

	// Convert amount to Wei
	ethAmount, success := new(big.Float).SetString(amount)
	if !success {
		return "", errors.New("invalid ETH amount")
	}
	weiAmount := new(big.Int)
	ethAmount.Mul(ethAmount, big.NewFloat(1e18)) // Convert ETH to Wei
	ethAmount.Int(weiAmount)

	// Get nonce (number of transactions sent from address)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	// Get current gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	// Estimate gas limit
	gasLimit := uint64(21000) // Standard gas limit for ETH transfer

	// Create transaction
	tx := types.NewTransaction(nonce, toAddress, weiAmount, gasLimit, gasPrice, nil)

	// Sign the transaction
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return "", err
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	// Return the transaction hash
	return signedTx.Hash().Hex(), nil
}

// FetchTransactionHistory fetches transactions for the current wallet
func FetchTransactionHistory() ([]Transaction, error) {
	// Ensure wallet exists
	err := EnsureWalletExists()
	if err != nil {
		return nil, err
	}

	// Get the wallet address
	address, err := GetAddress()
	if err != nil {
		return nil, err
	}

	// Get the Etherscan API key
	apiKey := os.Getenv("ETHERSCAN_API_KEY")
	if apiKey == "" {
		return nil, errors.New("Etherscan API key not set. Please set ETHERSCAN_API_KEY environment variable")
	}

	// Prepare the API request URL
	url := fmt.Sprintf(
		"https://api.etherscan.io/api?module=account&action=txlist&address=%s&sort=desc&apikey=%s",
		common.HexToAddress(address).Hex(), apiKey,
	)

	// Make the HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("failed to fetch transactions")
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response")
	}

	// Parse JSON response
	var result struct {
		Status  string        `json:"status"`
		Message string        `json:"message"`
		Result  []Transaction `json:"result"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New("failed to parse transaction history")
	}

	if result.Status != "1" {
		return nil, errors.New("failed to fetch transactions from Etherscan")
	}

	// Convert timestamps to readable format
	for i := range result.Result {
		ts, err := strconv.ParseInt(result.Result[i].Timestamp, 10, 64)
		if err == nil {
			result.Result[i].Timestamp = time.Unix(ts, 0).Format("2006-01-02 15:04:05")
		}
	}

	// Sort transactions by date (latest first)
	sort.Slice(result.Result, func(i, j int) bool {
		return result.Result[i].Timestamp > result.Result[j].Timestamp
	})

	return result.Result, nil
}

// GetPrivateKeyFromEnv retrieves the private key stored in the environment variable WALLET_PRIVATE_KEY
func GetPrivateKeyFromEnv() string {
	privateKey := os.Getenv(walletPrivateKey)
	return privateKey
}

// Ensure that a wallet exists before performing actions
func EnsureWalletExists() error {
	privateKey := GetPrivateKeyFromEnv()
	if privateKey == "" {
		return errors.New("wallet not found. Please create or import one")
	}
	return nil
}
