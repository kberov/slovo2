/*
Package cmd contains code for different actions (subcommands). Each file in
this package is an action of the application. We use `cobra` for managing
subcommands.

# Copyright © 2024-2025 Красимир Беров

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
	"errors"
	"os"

	"github.com/kberov/slovo2/slovo"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var Logger = slovo.Logger

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   slovo.Bin,
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
	Logger.Debug("in cmd.Execute")
	// Detect if we should execute the cgi command.
	if os.Getenv("GATEWAY_INTERFACE") != "" {
		Logger.Debug("in cmd.Execute GATEWAY_INTERFACE")
		os.Args = []string{os.Args[0], "cgi"}
		if err := cgiCmd.Execute(); err != nil {
			Logger.Error(err)
			os.Exit(1)
		}
		return
	}
	err := rootCmd.Execute()
	if err != nil {
		Logger.Error(err)
		os.Exit(1)
	}
}

var cfgFile = slovo.Cfg.File

func init() {
	cobra.EnableCommandSorting = false
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	pflags := rootCmd.PersistentFlags()
	pflags.StringVarP(&cfgFile, "config_file", "c", "",
		`Config file to use or you can set SLOVO_CONFIG environment variable to the
file to be read.  Alternatively we fall to sane internal defaults. See also
command 'config'.`)
	pflags.BoolVarP(&slovo.Cfg.Debug, "debug", "d", slovo.Cfg.Debug,
		"Display more verbose output in console.")
	// https://cobra.dev/#create-rootcmd
	// You will additionally define flags and handle configuration in your init() function.
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().StringVar("config_file", "c", , "Read configuration from this file")
	cobra.OnInitialize(rootInitConfig)
}

func rootInitConfig() {
	// If cfgFile is not passed on the command line and SLOVO_CONFIG is set,
	// use the path to config file in it.
	if len(cfgFile) == 0 {
		cfgFile = os.Getenv("SLOVO_CONFIG")
	}
	if len(cfgFile) == 0 {
		if slovo.Cfg.Debug {
			Logger.SetLevel(log.DEBUG)
		}
		return
	}
	// Try to Load YAML config if cfgFile exists. Otherwise
	// fallback to default configuration in slovo/config.go.
	finfo, err := os.Stat(cfgFile)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		Logger.Warnf("File %s does not exist. Falling back to internal configuration.", cfgFile)
	} else if finfo.Mode().IsRegular() && finfo.Mode().Perm()&0400 == 0400 {
		cfg, _ := os.ReadFile(cfgFile)
		if err := yaml.Unmarshal(cfg, &slovo.Cfg); err != nil {
			Logger.Fatal(err)
		} else {
			if !slovo.Cfg.Debug {
				Logger.SetLevel(log.INFO)
			}
		}
	}
}
