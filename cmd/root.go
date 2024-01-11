/*
Package cmd contains code for different actions. Each file is an action of the
application. We use `cobra` for managing subcommands. this is the root command.

# Copyright © 2024 Красимир Беров

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/kberov/slovo2/slovo"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

const defaultLogHeader = `${prefix}:${time_rfc3339}:${level}:${short_file}:${line}`

// Global koanf instance. Use . as the key path delimiter. This can be / or anything.
var (
	k      = koanf.New(".")
	parser = yaml.Parser()
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "slovo2",
	Short: "Наследникът на Слово – многократно по-бърз.",
	Long: `Наследникът на Слово – многократно по-бърз. Със запазен дух, но
изцяло осъществен наново на езика за програмиране Go. Автоматично открива и
работи в CGI среда.
ВНИМАНѤ!!!
Още сме в началото, така че има много грешки и недостатъци. Уча се!
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
	logger.Debug("in cmd.Execute")
	if os.Getenv("GATEWAY_INTERFACE") == "CGI/1.1" {
		logger.Debug("in cmd.Execute GATEWAY_INTERFACE")
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
var logger = log.New("slovo2")

func init() {
	logger.SetOutput(os.Stderr)
	logger.SetHeader(defaultLogHeader)
	cobra.EnableCommandSorting = false
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config_file", "c", cfgFile, "config file")
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
		logger.SetLevel(log.DEBUG)
	}
	//TODO: Load YAML config or use slovo.DefaultConfig.
	//if slovo.DefaultConfig.ConfigFile
	if err := k.Load(file.Provider(cfgFile), parser); err != nil {
		logger.Warnf("error loading config file: %v. Using slovo.DefaultConfig.", err)
	}
	logger.Debug("in root.go/rootInitConfig()")
}
