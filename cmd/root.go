package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "Wallet",
	Short: "Simple ETH wallet CLI",
	Long:  `This CLI allows you to create a wallet, check balance, and send ETH.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
}
