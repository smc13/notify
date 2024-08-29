package discord

import (
	"strings"
)

var discordURL = "https://discord.com/api/webhooks/"

// WebhookURLToParts takes a Discord webhook URL and returns the ID and token parts
func WebhookURLToParts(url string) (id, token string) {
	url = url[len(discordURL):]
	parts := strings.Split(url, "/")

	return parts[0], parts[1]
}
