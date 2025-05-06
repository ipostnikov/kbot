package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"gopkg.in/telebot.v3"
)

// downloads a video from Instagram using the provided URL
func DownloadInstagramVideo(url string) (string, error) {
	// define temp location
	tempFile := "/tmp/instagram_video.mp4"

	// run the yt-dlp command to download the vid in tempFile
	cmd := exec.Command("yt-dlp", "-o", tempFile, url)

	// Execute the command
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error downloading video with yt-dlp: %v", err)
	}

	// Return the path to the downloaded video file
	return tempFile, nil
}

// instagramCmd represents the instagram command
var instagramCmd = &cobra.Command{
	Use:   "instagram",
	Short: "Download Instagram videos using yt-dlp",
	Long:  `This command allows you to download Instagram videos by providing a URL.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Ensure that an Instagram URL is provided
		if len(args) < 1 {
			fmt.Println("Please provide an Instagram URL like: instagram <URL>")
			return
		}

		// get the URL from the arguments
		url := args[0]

		// Download the Instagram video
		tempFile, err := DownloadInstagramVideo(url)
		if err != nil {
			log.Printf("Error downloading video: %s", err)
			fmt.Println("Failed to download video.")
			return
		}

		// Initialize the Telegram bot using the existing Teletoken
		kbot, err := telebot.NewBot(telebot.Settings{
			Token: Teletoken, // Token defined  in kbot.go
		})
		if err != nil {
			log.Fatalf("Error creating bot: %v", err)
			return
		}

		// Handle text messages to download the video
		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			// Check if the message is the correct command to trigger the download
			if m.Text() == "/download" {
				// Send the downloaded video as a file to the user
				err := m.Send(&telebot.Video{
					File: telebot.FromDisk(tempFile), // Use FromDisk to send the video file
				})
				if err != nil {
					log.Printf("Failed to send video: %v", err)
					return err
				}

				// After sending the video, clean up the temporary file
				defer os.Remove(tempFile) // Ensure that the file is removed once sent

				return nil
			}

			return nil
		})

		// Start the bot to listen for messages
		kbot.Start()
	},
}

func init() {
	// Add the instagram command to the root command
	rootCmd.AddCommand(instagramCmd)
}
