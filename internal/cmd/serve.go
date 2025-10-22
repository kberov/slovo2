package cmd

import (
	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command.
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run Slovo as a Sever.",
	Long:  `Starts Slovo as a HTTP server.`,
	// I had to move init* functions here to make sure that only the parent and
	// respective command's init* are run.
	PreRun: func(_ *cobra.Command, args []string) {
		serveInitConfig()
		Logger.Debugf("serveCmd.Command().PreRun() called. args: %+v ", args)
	},
	Run: func(_ *cobra.Command, args []string) {
		Logger.Debugf("serveCmd.Command().Run() called with args %+v", args)
		slovo.Start(Logger)
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
	serveCmd.Flags().StringVarP(&slovo.Cfg.Serve.Location, "listen", "l",
		slovo.Cfg.Serve.Location, "Location to listen on")
	// cobra.OnInitialize(serveInitConfig)
}

func serveInitConfig() {
	Logger.Debug("in serve.go/serveInitConfig()")
	Logger.Debugf("Listening on %s.", slovo.Cfg.Serve.Location)
}
