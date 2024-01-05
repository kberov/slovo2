/*
Copyright © 2024 Красимир Беров
*/
package cmd

import (
	"fmt"

	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run Slovo as a Sever.",
	Long:  `Starts Slovo as a HTTP server.`,
	Run: func(cmd *cobra.Command, args []string) {
		slovo.Serve()
		fmt.Println("serve called")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
