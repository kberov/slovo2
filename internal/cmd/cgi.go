package cmd

import (
	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
)

// cgiCmd represents the cgi command.
var cgiCmd = &cobra.Command{
	Use:   "cgi",
	Short: spf("Run %s as a CGI script.", slovo.Bin),
	Long: spf(`
This command will be executed automatically by %[1]s if the GATEWAY_INTERFACE
environment variable is set. In other words %[1]s autodetects from the
environment, if it is invoked by a web server like Apache or as a commandline
application. Also this is how we cheat %[1]s to test it on the command line.
`, slovo.Bin),
	PreRun: cgiPreRun,
	Run: func(_ *cobra.Command, args []string) {
		_ = args
		slovo.StartCGI(Logger)
	},
}

func init() {
	rootCmd.AddCommand(cgiCmd)
	slovo.CgiInitFlags(cgiCmd.Flags())
}

func cgiPreRun(_ *cobra.Command, args []string) {
	Logger.Debugf("cgiCmd.PreRun(cgiCmd): args: %v", args)
}
