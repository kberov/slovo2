package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var spf = fmt.Sprintf
var defaultCfg = slovo.Cfg

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config [action]",
	Short: "A command to manage slovo2 configuration",
	Long: spf(`This command performs various actions with the configuration.
Available actions are:
  defaults - Displays the default configuration.
  dump     - Dumps the configuration to specified file with --config_file.
             Defults to %s. 
`, defaultCfg.ConfigFile),
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		switch args[0] {
		case `defaults`:
			displayDefaultConfig()
		case `dump`:
			dumpConfig()
		default:
			fmt.Println("\nUnknown action!")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func displayDefaultConfig() {
	cfg, err := yaml.Marshal(&defaultCfg)
	if err != nil {
		Logger.Fatalf("error: %v", err)
	}
	fmt.Printf("---\n%s\n\n", string(cfg))
}

func dumpConfig() {
	if cfgFile == "" {
		cfgFile = defaultCfg.ConfigFile
	}
	fmt.Printf("Will dump configuration to %s.\n\n\tNote! The directory must exist.\n\n", cfgFile)
	fmt.Printf("Default configuration file is %s.\n", defaultCfg.ConfigFile)
	finfo, err := os.Stat(cfgFile)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// fine
	} else if finfo.Mode().IsRegular() && finfo.Mode().Perm()&0400 == 0400 {
		// backup the existing file if it exists
		if e := os.Rename(cfgFile, cfgFile+".old"); e != nil {
			Logger.Fatal(e)
		}
	}
	cfg, err := yaml.Marshal(&defaultCfg)
	if err != nil {
		Logger.Fatalf("error: %v", err)
	}
	if err := os.WriteFile(cfgFile, cfg, 0600); err != nil {
		Logger.Fatalf(`%s: %s`, cfgFile, err.Error())
	} else {
		fmt.Printf("Configuration dumped to %s\n", cfgFile)
	}
}
