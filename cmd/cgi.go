/*
Copyright © 2024 Красимир Беров
*/
package cmd

import (
	"fmt"

	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
)

// cgiCmd represents the cgi command
var cgiCmd = &cobra.Command{
	Use:   "cgi",
	Short: "Run Slovo as a CGI script.",
	Long: `This command will be executed automatically if the GATEWAY_INTERFACE
environment variable is set.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cgiCmd.PreRun(cgiCmd): args: %v\n", args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		slovo.ServeCGI()
		//panic("how we got here?")
		//fmt.Printf("cgiCmd.Run(cmd) cmd: %#v\n", cmd)
	},
}

func init() {
	rootCmd.AddCommand(cgiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cgiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cgiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
