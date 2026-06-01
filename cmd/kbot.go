package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v3"
)

var (
	appVersion = "dev"
	Teletoken  = os.Getenv("TELE_TOKEN")
	// Optional path to a Netscape-format cookies file for Instagram.
	// Many reels require a logged-in session; set this to a mounted cookies.txt.
	cookiesFile = os.Getenv("INSTAGRAM_COOKIES")
)

var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Start the Telegram bot",
	Long:    `Start the Telegram bot. Send an Instagram URL to download and receive the video.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kbot %s started\n", appVersion)

		kbot, err := telebot.NewBot(telebot.Settings{
			Token:  Teletoken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})
		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
		}

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {
			messageText := m.Text()

			if messageText == "/start" {
				return m.Send(fmt.Sprintf("Hello! I'm K9 Shepard Bot %s.\nSend me an Instagram video URL to download.", appVersion))
			}

			if strings.Contains(messageText, "instagram.com") {
				log.Printf("Instagram URL received: %s", messageText)

				if err := m.Send("Downloading your Instagram video, please wait..."); err != nil {
					log.Printf("Failed to notify user: %v", err)
				}

				tempFile, err := downloadInstagramVideo(messageText)
				if err != nil {
					log.Printf("Error downloading video: %s", err)
					return m.Send("Failed to download the video.")
				}
				defer os.RemoveAll(filepath.Dir(tempFile))

				video := &telebot.Video{File: telebot.FromDisk(tempFile)}
				if err = m.Send(video); err != nil {
					log.Printf("Failed to send video: %v", err)
					return m.Send("Failed to send the video.")
				}
			}

			return nil
		})

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-stop
			kbot.Stop()
		}()

		kbot.Start()
	},
}

// downloadInstagramVideo downloads url into a fresh temp directory and returns
// the path to the resulting file. The caller is responsible for removing the
// file's parent directory (e.g. os.RemoveAll(filepath.Dir(path))).
func downloadInstagramVideo(url string) (string, error) {
	// Download into a temp directory rather than a pre-created file: yt-dlp
	// skips the download ("already downloaded") if the output path exists, so
	// the destination must not exist yet.
	dir, err := os.MkdirTemp("", "kbot-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %v", err)
	}
	out := filepath.Join(dir, "video.mp4")

	// Request a single pre-muxed file so yt-dlp never has to merge separate
	// audio/video streams — that lets us ship the image without ffmpeg.
	args := []string{
		"-f", "best[ext=mp4]/best",
		"--no-playlist",
		"-o", out,
	}
	if cookiesFile != "" {
		args = append(args, "--cookies", cookiesFile)
	}
	args = append(args, url)

	cmd := exec.Command("yt-dlp", args...)
	if combined, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(dir)
		return "", fmt.Errorf("yt-dlp: %v: %s", err, combined)
	}

	// Guard against an empty result (yt-dlp can exit 0 without writing a file);
	// otherwise Telegram rejects the upload with "file must be non-empty".
	if fi, err := os.Stat(out); err != nil || fi.Size() == 0 {
		os.RemoveAll(dir)
		return "", fmt.Errorf("yt-dlp produced no video file")
	}

	return out, nil
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}
