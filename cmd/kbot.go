package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	appVersion = "1.0.3" // Define your application version here
	Teletoken  = os.Getenv("TELE_TOKEN")
)

// kbotCmd represents the kbot command
var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Telegram bot for interacting with Instagram videos",
	Long:    `A bot to download Instagram videos by providing a URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started\n", appVersion)

		kbot, err := telebot.NewBot(telebot.Settings{
			URL:    "",
			Token:  Teletoken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
			return
		}

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			messageText := m.Text()

			if messageText == "/start" {
				return m.Send(fmt.Sprintf("Hello! I'm Kbot %s.\nSend me an Instagram video URL to download.", appVersion))
			}

			if strings.Contains(messageText, "instagram.com") {
				log.Printf("Instagram URL received: %s", messageText)

				err := m.Send("Downloading your Instagram video, please wait...")
				if err != nil {
					log.Printf("Failed to notify user: %v", err)
				}

				tempFile, err := DownloadInstagramVideo(messageText)
				if err != nil {
					log.Printf("Error downloading video: %s", err)
					return m.Send("❌ Failed to download the video.")
				}
				defer os.Remove(tempFile)

				video := &telebot.Video{File: telebot.FromDisk(tempFile)}
				err = m.Send(video)
				if err != nil {
					log.Printf("Failed to send video: %v", err)
					return m.Send("❌ Failed to send the video.")
				}
			}

			return nil
		})

		kbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}
