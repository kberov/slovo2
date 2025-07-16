//lint:file-ignore ST1000 Already documented in root.go

/*
Copyright © 2024 Красимир Беров
*/
package cmd

import (
	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
)

// cgiCmd represents the cgi command
var cgiCmd = &cobra.Command{
	Use:   "cgi",
	Short: "Run Slovo as a CGI script.",
	Long: `This command will be executed automatically if the GATEWAY_INTERFACE
environment variable is set. In other words Slovo2 autodetects from the
environment, if it is invoked by a web server like Apache or as a commandline
application. Also this is how we cheat Slovo2 to test it on the command line.`,
	// I had to move init* functions here to make sure that only the parent and
	// respective command's init* are run.
	PreRun: func(cmd *cobra.Command, args []string) {
		// Logger.Debugf("cgiCmd.PreRun(cgiCmd): args: %v", args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		slovo.StartCGI(Logger)
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

	slovo.CgiInitFlags(cgiCmd.Flags())
}
