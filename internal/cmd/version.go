//lint:file-ignore ST1000 Already documented in root.go
/*
Copyright © 2024 Красимир Беров
*/
package cmd

import (
	"fmt"

	"github.com/kberov/slovo2/slovo"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Текущо издание на Слово",
	Long: `Пълно описание на изданието на Слово. Показва дата на изданието и
кода. Кодът е буква от глаголицата и обозначава някаква стъпка в развитието
на приложението.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(
			`VERSION: %s
CODENAME: %s 

`, slovo.VERSION, slovo.CODENAME)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
