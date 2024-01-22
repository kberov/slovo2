//lint:file-ignore ST1000 Already documented in root.go

/*
Copyright © 2024 Красимир Беров
*/
package cmd

import (
	"os"

	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
)

// cgiCmd represents the cgi command
var cgiCmd = &cobra.Command{
	Use:   "cgi",
	Short: "Run Slovo as a CGI script.",
	Long: `This command will be executed automatically if the GATEWAY_INTERFACE
environment variable is set.`,
	// I had to move init* functions here to make sure that only the parent and
	// respective command's init* are run.
	PreRun: func(cmd *cobra.Command, args []string) {
		cgiInitConfig()
		Logger.Debugf("cgiCmd.PreRun(cgiCmd): args: %v", args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		slovo.ServeCGI(Logger)
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
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.ServeCGI.HTTP_HOST,
		"HTTP_HOST", "H",
		slovo.Cfg.ServeCGI.HTTP_HOST, "The server host to which the client request is directed.")
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.ServeCGI.REQUEST_URI,
		"REQUEST_URI", "U",
		slovo.Cfg.ServeCGI.REQUEST_URI, "Request URI")
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.ServeCGI.REQUEST_METHOD,
		"REQUEST_METHOD", "M",
		slovo.Cfg.ServeCGI.REQUEST_METHOD, "Request method")
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.ServeCGI.SERVER_PROTOCOL,
		"SERVER_PROTOCOL", "P",
		slovo.Cfg.ServeCGI.SERVER_PROTOCOL, "Server protocol")
	//cobra.OnInitialize(rootInitConfig)
	//cobra.OnInitialize(cgiInitConfig)
}

func cgiInitConfig() {
	// TODO
	Logger.Debugf("in cgiCmd.cgiInitConfig()")

	// minimum ENV values for emulating a CGI request on the command line
	var env = map[string]string{
		"SERVER_PROTOCOL": slovo.Cfg.ServeCGI.SERVER_PROTOCOL,
		"REQUEST_METHOD":  slovo.Cfg.ServeCGI.REQUEST_METHOD,
		"HTTP_HOST":       slovo.Cfg.ServeCGI.HTTP_HOST,
		//"HTTP_REFERER":        "elsewhere",
		//"HTTP_USER_AGENT":     "slovo2client",
		"HTTP_ACCEPT_CHARSET": "utf-8",
		// "HTTP_FOO_BAR":    "baz",
		"REQUEST_URI": slovo.Cfg.ServeCGI.REQUEST_URI,
		// "CONTENT_LENGTH":  "123",
		"CONTENT_TYPE": "text/html",
		// "REMOTE_ADDR":     "5.6.7.8",
		// "REMOTE_PORT":     "54321",
	}
	for k, v := range env {
		if os.Getenv(k) == "" {
			Logger.Debugf("Setting %s: %s", k, v)
			os.Setenv(k, v)
		}
	}
}
