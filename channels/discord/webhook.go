package discord

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/smc13/notify"
)

type DiscordWebhookChannel struct {
	webhookID    string
	webhookToken string
	session      *discordgo.Session
	params       discordgo.WebhookParams
}

type DiscordWebhookNotification interface {
	ToDiscordWebhook(notify.Notifiable) (*discordgo.WebhookParams, error)
}

type DiscordWebhookNotifiable interface {
	DiscordUserIds() []string
}

func NewWebhookChannel(id string, token string) *DiscordWebhookChannel {
	session, err := discordgo.New("")
	// we don't actually need a bot here since we're using webhooks
	// just panic for now if something changes in discordgo
	if err != nil {
		panic(err)
	}

	return &DiscordWebhookChannel{
		webhookID:    id,
		webhookToken: token,
		session:      session,
	}
}

// SetParams sets the webhook params to be sent to the webhook when not using a DiscordWebhookNotification
func (d *DiscordWebhookChannel) SetParams(params discordgo.WebhookParams) {
	d.params = params
}

func (d *DiscordWebhookChannel) Notify(ctx context.Context, notifiable notify.Notifiable, notif notify.Notification) error {
	message, err := d.notificationToDiscordEmbed(notifiable, notif)
	if err != nil {
		return err
	}

	if message == nil {
		return nil
	}

	if notifiable, ok := notifiable.(DiscordWebhookNotifiable); ok {
		userIds := notifiable.DiscordUserIds()
		formattedChannels := make([]string, len(userIds))
		for i, userId := range userIds {
			formattedChannels[i] = fmt.Sprintf("<@%s>", userId)
		}

		message.Content = fmt.Sprintf("%s\n%s", message.Content, strings.Join(formattedChannels, " "))
	}

	_, err = d.session.WebhookExecute(d.webhookID, d.webhookToken, true, message, discordgo.WithContext(ctx))
	return err
}

func (d *DiscordWebhookChannel) notificationToDiscordEmbed(notifiable notify.Notifiable, notif notify.Notification) (*discordgo.WebhookParams, error) {
	if discordNotif, ok := notif.(DiscordWebhookNotification); ok {
		return discordNotif.ToDiscordWebhook(notifiable)
	}

	params := d.params
	params.Embeds = []*discordgo.MessageEmbed{
		{Title: notif.Subject(), Description: notif.Content()},
	}

	return &params, nil
}
