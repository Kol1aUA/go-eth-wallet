package cmd

import (
	"fmt"
	"go-eth-wallet/internal/wallet"
	"log"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate a new ETH wallet",
	Run: func(cmd *cobra.Command, args []string) {
		address, err := wallet.CreateWallet()
		if err != nil {
			log.Fatal("Error creating wallet:", err)
		}
		fmt.Println("Wallet created successfully!")
		fmt.Println("Address:", address)
	},
}

var importComand = &cobra.Command{
	Use:   "import <private-key>",
	Short: "Import ETH wallet using private key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		privateKey := args[0]
		err := wallet.ImportWallet(privateKey)
		if err != nil {
			log.Fatal("Error importing wallet:", err)
		}

		fmt.Println("Wallet imported successfully!")
	},
}

var showCommand = &cobra.Command{
	Use:   "show",
	Short: "Show the current wallet`s ETH address",
	Run: func(cmd *cobra.Command, args []string) {
		address, err := wallet.GetAddress()
		if err != nil {
			log.Fatal("Error fetching wallet address:", err)
		}

		fmt.Println("Wallet address:", address)
	},
}

var balanceCommand = &cobra.Command{
	Use:   "balance",
	Short: "Check ETH balance of the current wallet",
	Run: func(cmd *cobra.Command, args []string) {
		balance, err := wallet.GetBalance()
		if err != nil {
			log.Fatal("Error fetching balance:", err)
		}
		fmt.Println("Wallet Balance:", balance, "ETH")
	},
}

var sendCommand = &cobra.Command{
	Use:   "send <recipient-address> <amount>",
	Short: "Send ETH to another Ethereum address",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		recipient := args[0]
		amount := args[1]

		txHash, err := wallet.SendETH(recipient, amount)
		if err != nil {
			log.Fatal("Error sending ETH:", err)
		}

		fmt.Println("Transaction sent successfully!")
		fmt.Println("Transaction Hash:", txHash)
	},
}

var historyCommand = &cobra.Command{
	Use:   "history",
	Short: "Fetch and display Ethereum transaction history",
	Run: func(cmd *cobra.Command, args []string) {
		transactions, err := wallet.FetchTransactionHistory()
		if err != nil {
			log.Fatal("Error fetching transaction history:", err)
		}

		fmt.Println("Transaction History:")
		fmt.Println("-------------------------------------------------")
		for _, tx := range transactions {
			fmt.Printf("🆔 Tx Hash: %s\n", tx.Hash)
			fmt.Printf("🔄 From: %s\n", tx.From)
			fmt.Printf("➡️  To: %s\n", tx.To)
			fmt.Printf("💰 Value: %s ETH\n", tx.Value)
			fmt.Printf("📅 Date: %s\n", tx.Timestamp)
			fmt.Println("-------------------------------------------------")
		}
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
	RootCmd.AddCommand(importComand)
	RootCmd.AddCommand(showCommand)
	RootCmd.AddCommand(balanceCommand)
	RootCmd.AddCommand(sendCommand)
	RootCmd.AddCommand(historyCommand)
}
