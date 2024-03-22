package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"go-code/awesomeProject1/app/router"
	"os"
)

var port int

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "A brief description of your application",
	Run: func(cmd *cobra.Command, args []string) {
		// Start the server
		router.New(port)
	},
}

func init() {
	// Define the port flag
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port number to run the server on")
}
func main() {
	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
