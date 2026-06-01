# 🐾 k9_shepard_bot

A small Telegram bot that downloads Instagram videos. Send it a reel or
post URL and it sends the video back. Uses [`yt-dlp`](https://github.com/yt-dlp/yt-dlp)
to do the actual downloading.

Try it: [@k9_shepard_bot](https://t.me/k9_shepard_bot)

<p align="center">
<img src="usage.gif" alt="Usage example" width="600">
</p>

## Run it

You need Docker and a bot token from [@BotFather](https://t.me/BotFather).

```bash
git clone https://github.com/ipostnikov/kbot.git
cd kbot

echo "TELE_TOKEN=your-token-here" > .env
docker compose up -d
```

That's it. Logs: `docker compose logs -f`. Stop it: `docker compose down`.

## Instagram cookies (optional)

Some posts require a logged-in session. If a download fails with
"Instagram API is not granting access", give it your cookies:

1. Export your Instagram cookies to `cookies.txt` (Netscape format) — the
   "Get cookies.txt LOCALLY" browser extension works well.
2. Put `cookies.txt` next to `docker-compose.yml`.
3. In `docker-compose.yml`, uncomment the `volumes` block and add
   `INSTAGRAM_COOKIES=/cookies.txt` to your `.env`.
4. `docker compose up -d`.

## License

MIT — see [LICENSE](LICENSE).
