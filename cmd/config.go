package cmd

import (
	"fmt"

	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A command to dump to file the default configuration",
	Long:  `This command dumps he default configuration structure slovo.Cfg to a yaml file using yaml.v3`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
		dumpConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dumpConfig() {

	d, err := yaml.Marshal(&slovo.Cfg.DB)
	if err != nil {
		Logger.Fatalf("error: %v", err)
	}
	Logger.Debugf("--- DBConfig to yaml:\n%s\n\n", string(d))
	// Back to struct
	var db slovo.DBConfig
	if err := yaml.Unmarshal(d, &db); err != nil {
		Logger.Fatalf("%w", err)
	}
	Logger.Debugf("Back to slovo.DBConfig: %#v ", db)

	Logger.Debug(`--------------`)

	cgikfg, err := yaml.Marshal(&slovo.Cfg.ServeCGI)

	if err != nil {
		Logger.Fatalf("error: %v", err)
	}

	Logger.Debugf("--- cgikfg to yaml:\n%s\n\n", string(cgikfg))

	// Back to struct
	var cgistruct slovo.ServeCGIConfig
	if err := yaml.Unmarshal(cgikfg, &cgistruct); err != nil {
		Logger.Fatalf("%w", err)
	}
	Logger.Debugf("Back to slovo.ServeCGIConfig: %#v ", cgistruct)

	Logger.Debug(`--------------`)

	routes, err := yaml.Marshal(&slovo.Cfg.Routes)

	if err != nil {
		Logger.Fatalf("error: %v", err)
	}

	Logger.Debugf("--- routes to yaml:\n%s\n\n", string(routes))

	// Back to struct
	var routesslice slovo.Routes
	if err := yaml.Unmarshal(routes, &routesslice); err != nil {
		Logger.Fatalf("%w", err)
	}
	Logger.Debugf("Back to slovo.Routes: %#v ", routesslice)

	Logger.Debug(`--------------`)

	rewrites, err := yaml.Marshal(&slovo.Cfg.RewriteConfig)

	if err != nil {
		Logger.Fatalf("error: %v", err)
	}

	Logger.Debugf("--- rewrites to yaml:\n%s\n\n", string(rewrites))

	// Back to struct
	var rewriterules slovo.RewriteConfig
	if err := yaml.Unmarshal(rewrites, &rewriterules); err != nil {
		Logger.Fatalf("%w", err)
	}
	Logger.Debugf("Back to slovo.Routes: %#v ", rewriterules)

	Logger.Debug(`ALLL--------------`)

	cfg, err := yaml.Marshal(&slovo.Cfg)

	if err != nil {
		Logger.Fatalf("error: %v", err)
	}

	Logger.Debugf("--- cfg to yaml:\n%s\n\n", string(cfg))

	// Back to struct
	var cfgStruct slovo.Config
	if err := yaml.Unmarshal(cfg, &cfgStruct); err != nil {
		Logger.Fatalf("%w", err)
	}
	Logger.Debugf("Back to slovo.Config: %#v ", cfgStruct)
}
