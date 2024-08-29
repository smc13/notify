package discord

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebhookURLToParts(t *testing.T) {
	// I've broken the url... so don't bother trying to use it
	url := "https://discord.com/api/webhooks/115ffff895182114916/E7Wm7V0v9xf3f3fzM1kbsEOHF2DPgmsNq6f6A4XUEDQFV1EIH0ZsRH9S68RfsfsfUcIM"
	id, token := WebhookURLToParts(url)
	assert.Equal(t, "115ffff895182114916", id)
	assert.Equal(t, "E7Wm7V0v9xf3f3fzM1kbsEOHF2DPgmsNq6f6A4XUEDQFV1EIH0ZsRH9S68RfsfsfUcIM", token)
}
