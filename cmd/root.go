package cmd

import (
	"fmt"
	"os"

	"github.com/kberov/slovo2/slovo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

// const defaultLogHeader = `${prefix}:${time_rfc3339}:${level}:${short_file}:${line}`
// See possible placeholders in  github.com/labstack/gommon@v0.4.2/log/log.go function log()
// See https://echo.labstack.com/docs/customization
const defaultLogHeader = `${prefix}:${level}:${short_file}:${line}`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slovo2",
	Short: "Наследникът на Слово – многократно по-бърз.",
	Long: `Наследникът на Слово – многократно по-бърз. Със запазен дух, но
изцяло осъществен наново на езика за програмиране Go. Автоматично открива и
работи в CGI среда.
ВНИМАНѤ!!! Още сме в началото, така че има много грешки и недостатъци. Уча се!
`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	/*
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("rootCmd.Run(rootCmd): args: %v\n", args)
		},
	*/
	//	PreRun: func(cmd *cobra.Command, args []string) {
	//		rootInitConfig()
	//		//fmt.Printf("rootCmd.Run(rootCmd): args: %v\n", args)
	//	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Logger.Debug("in cmd.Execute")
	if os.Getenv("GATEWAY_INTERFACE") == "CGI/1.1" {
		// Logger.Debug("in cmd.Execute GATEWAY_INTERFACE")
		os.Args = []string{os.Args[0], "cgi"}
		if cgiCmd.Execute() != nil {
			os.Exit(1)
		}
		return
	}
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile = slovo.Cfg.ConfigFile
var Logger = log.New("slovo2")

func init() {
	Logger.SetOutput(os.Stderr)
	Logger.SetHeader(defaultLogHeader)

	cobra.EnableCommandSorting = false
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config_file", "c", slovo.Cfg.ConfigFile, "config file")
	rootCmd.PersistentFlags().BoolVarP(&slovo.Cfg.Debug, "debug", "d", slovo.Cfg.Debug,
		"Display more verbose output in console output. default: "+fmt.Sprintf("%v", slovo.Cfg.Debug))
	// https://cobra.dev/#create-rootcmd
	// You will additionally define flags and handle configuration in your init() function.
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().StringVar("config_file", "c", , "Read configuration from this file")
	cobra.OnInitialize(rootInitConfig)
}

func rootInitConfig() {
	if slovo.Cfg.Debug {
		Logger.SetLevel(log.DEBUG)
	}
	// TODO: Try to Load YAML config if cfgFile != slovo.Cfg.ConfigFile - default value
	// if slovo.Cfg.ConfigFile....
	// if cfgFile != slovo.Cfg.ConfigFile {
	// TODO:
	// }
	// Logger.Debug("in root.go/rootInitConfig()")
}
