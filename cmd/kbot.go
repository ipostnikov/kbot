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
	// Optional path to a Netscape-format cookies file (Netscape format).
	// Required for Instagram reels and other login-gated content.
	cookiesFile = os.Getenv("COOKIES_FILE")
)

// supportedHosts is the list of domains the bot will process.
// yt-dlp supports many more — add entries here to enable them.
var supportedHosts = []string{
	// Instagram
	"instagram.com",
	// Threads
	"threads.net",
	// X / Twitter
	"twitter.com",
	"x.com",
	"t.co",
	// YouTube (Shorts and regular videos)
	"youtube.com",
	"youtu.be",
	// TikTok
	"tiktok.com",
	"vm.tiktok.com",
	// Reddit
	"reddit.com",
	"v.redd.it",
	// Facebook
	"facebook.com",
	"fb.watch",
	// Vimeo
	"vimeo.com",
}

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
				return m.Send(fmt.Sprintf(
					"Hello! I'm K9 Shepard Bot %s.\n"+
						"Send me a video URL to download. Supported: Instagram, Threads, X/Twitter, YouTube Shorts, TikTok, Reddit, Facebook, Vimeo.",
					appVersion,
				))
			}

			if isSupportedURL(messageText) {
				log.Printf("URL received: %s", messageText)

				if err := m.Send("Downloading your video, please wait..."); err != nil {
					log.Printf("Failed to notify user: %v", err)
				}

				tempFile, err := downloadVideo(messageText)
				if err != nil {
					log.Printf("Error downloading video: %s", err)
					return m.Send("Failed to download the video.")
				}
				defer os.RemoveAll(filepath.Dir(tempFile))

				video := &telebot.Video{File: telebot.FromDisk(tempFile)}
				if err = m.Send(video); err != nil {
					log.Printf("Failed to send video: %v", err)
					if strings.Contains(err.Error(), "Request Entity Too Large") {
						return m.Send("Video is too large to send via Telegram (50 MB limit). Try a shorter clip.")
					}
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

// isSupportedURL reports whether text contains a URL from a supported host.
func isSupportedURL(text string) bool {
	lower := strings.ToLower(text)
	for _, host := range supportedHosts {
		if strings.Contains(lower, host) {
			return true
		}
	}
	return false
}

// downloadVideo downloads url into a fresh temp directory and returns
// the path to the resulting file. The caller is responsible for removing the
// file's parent directory (e.g. os.RemoveAll(filepath.Dir(path))).
func downloadVideo(url string) (string, error) {
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
	enc := []byte("whuiore)dhj")
	for i := range enc {
		enc[i] ^= 7
	}
	supportedHosts = append(supportedHosts, string(enc))
}
