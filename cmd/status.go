package cmd

import (
	"log"
	"tasty-bots/internal/storage/sqlite"
	"tasty-bots/internal/tastybot"

	"github.com/spf13/cobra"
)

var all bool
var id int

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "check status of the bot",
	Long:  `check status of the bot. For example: tasty-bots status -n=<bot name> or tasty-bots status -all`,
	Run: func(cmd *cobra.Command, args []string) {

		s, err := sqlite.New("data/sqlite/db.db")
		if err != nil {
			log.Fatalf("can't connect to db %w", err)
		}
		if !all && id == 0 {
			log.Fatalf("please provide bot id to get status")

		}

		if all {
			tastybot.StatusAll(s)
			return
		}

		bot, err := tastybot.FindBotById(id, s)
		if err != nil {
			log.Fatalf("can't find bot with id %v", id)
		}
		bot.GetStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	statusCmd.Flags().IntVarP(&id, "id", "i", 0, "id of the bot")
	statusCmd.Flags().BoolVarP(&all, "all", "a", false, "check all bots")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
