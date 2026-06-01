package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kbot",
	Short: "Telegram bot for downloading Instagram videos",
	Long:  `kbot is a Telegram bot that downloads Instagram videos via yt-dlp and sends them back to the user.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
