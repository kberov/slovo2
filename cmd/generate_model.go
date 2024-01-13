package cmd

import (
	"fmt"

	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
)

// generate/modelCmd represents the generate/model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generates Go code for sqlx",
	Long: `
This command generates idiomatic (I'm trying) Go code for using sqlx with your
own queries and data objects almost like an ORM. We use sqlite3 and extract all
the meta data for the model from the database. Please, pass the path to your
database file or add it to the configuration section Cfg.Db.DSN
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate/model called")
	},
}

func init() {
	generateCmd.AddCommand(modelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generate/modelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generate/modelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	modelCmd.Flags().StringVarP(&slovo.Cfg.DB.DSN, "DSN", "D", slovo.Cfg.DB.DSN, "DSN for the database")
	modelCmd.Flags().StringVarP(&slovo.Cfg.DB.tables, "tables", "t", slovo.Cfg.DB.tables, "tables for which to generate model types")
}
