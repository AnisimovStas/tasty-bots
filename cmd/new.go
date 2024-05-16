package cmd

import (
	"log"
	"tasty-bots/internal/storage/sqlite"
	"tasty-bots/internal/tastybot"

	"github.com/spf13/cobra"
)

const (
	DEFAULT_URL = "https://tastydrop.in"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Command to create new bot",
	Long:  `Command to create new bot For example: tasty-bots new -t=<tasty token>`,
	Run: func(cmd *cobra.Command, args []string) {
		t := cmd.Flag("token").Value.String()
		u := cmd.Flag("url").Value.String()
		if u == "" {
			u = DEFAULT_URL
		}
		s, err := sqlite.New("data/sqlite/db.db")
		if err != nil {
			log.Fatalf("can't connect to db %w", err)
		}

		tastybot.New(t, u, s)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("token", "t", "", "tasty token")
	newCmd.Flags().StringP("url", "u", "", "base url")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
