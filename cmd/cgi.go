//lint:file-ignore ST1000 Already documented in root.go

/*
Copyright © 2024 Красимир Беров
*/
package cmd

import (
	"net/url"
	"os"
	"strings"

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
		cgiInitConfig()
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
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.StartCGI.HTTP_HOST,
		"HTTP_HOST", "H",
		slovo.Cfg.StartCGI.HTTP_HOST, "The server host to which the client request is directed.")
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.StartCGI.REQUEST_URI,
		"REQUEST_URI", "U",
		slovo.Cfg.StartCGI.REQUEST_URI, "Request URI")
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.StartCGI.REQUEST_METHOD,
		"REQUEST_METHOD", "M",
		slovo.Cfg.StartCGI.REQUEST_METHOD, "Request method")
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.StartCGI.SERVER_PROTOCOL,
		"SERVER_PROTOCOL", "P",
		slovo.Cfg.StartCGI.SERVER_PROTOCOL, "Server protocol")
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.StartCGI.HTTP_ACCEPT_CHARSET,
		"HTTP_ACCEPT_CHARSET", "C",
		slovo.Cfg.StartCGI.HTTP_ACCEPT_CHARSET, "Accept-Charset")
	cgiCmd.Flags().StringVarP(
		&slovo.Cfg.StartCGI.CONTENT_TYPE,
		"CONTENT_TYPE", "T",
		slovo.Cfg.StartCGI.CONTENT_TYPE, "Content-Type")
}

func cgiInitConfig() {
	if os.Getenv("GATEWAY_INTERFACE") != "" {
		return
	}

	// minimum ENV values for emulating a CGI request on the command line
	var env = map[string]string{
		"GATEWAY_INTERFACE": "CGI/1.1",
		"SERVER_PROTOCOL":   slovo.Cfg.StartCGI.SERVER_PROTOCOL,
		"REQUEST_METHOD":    slovo.Cfg.StartCGI.REQUEST_METHOD,
		"HTTP_HOST":         slovo.Cfg.StartCGI.HTTP_HOST,
		//"HTTP_REFERER":        "elsewhere",
		//"HTTP_USER_AGENT":     "slovo2client",
		"HTTP_ACCEPT_CHARSET": slovo.Cfg.StartCGI.HTTP_ACCEPT_CHARSET,
		// "HTTP_FOO_BAR":    "baz",
		"REQUEST_URI": escapeRequestURI(slovo.Cfg.StartCGI.REQUEST_URI),
		// "CONTENT_LENGTH":  "123",
		"CONTENT_TYPE": slovo.Cfg.StartCGI.CONTENT_TYPE,
		// "REMOTE_ADDR":     "5.6.7.8",
		// "REMOTE_PORT":     "54321",
	}
	for k, v := range env {
		if os.Getenv(k) == "" {
			// Logger.Debugf("Setting %s: %s", k, v)
			os.Setenv(k, v)
		}
	}
}

// escapeRequestURI escapes the passed on the commandline address as if a user
// agent did it.
func escapeRequestURI(uri string) string {
	uri, _ = strings.CutPrefix(uri, `/`)
	parts := strings.Split(uri, `/`)
	uri = ""
	for _, p := range parts {
		uri += `/` + url.PathEscape(p)
	}
	return uri
}
