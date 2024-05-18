/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"tasty-bots/internal/storage/sqlite"
	"tasty-bots/internal/tastybot"

	"github.com/spf13/cobra"
)

var caseTD string

// runCasesCmd represents the runCases command
var runCasesCmd = &cobra.Command{
	Use:   "runCases",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		s, err := sqlite.New("data/sqlite/db.db")
		if err != nil {
			log.Fatalf("can't connect to db %w", err)
		}

		err = tastybot.RunCases(id, caseTD, s)
		if err != nil {
			log.Fatalf("error on set status %w", err)

		}

	},
}

func init() {
	rootCmd.AddCommand(runCasesCmd)
	runCasesCmd.Flags().IntVarP(&id, "id", "i", 0, "id of the bot")
	runCasesCmd.Flags().StringVarP(&caseTD, "case", "c", "v2_rare", "id of the bot")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCasesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCasesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
