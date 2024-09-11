package ntfy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/smc13/notify"
)

type NtfyChannel struct {
	baseURL string
	topic   string
	client  *http.Client
}

type NtfyNotification interface {
	ToNtfy(notify.Notifiable) (Message, error)
}

type NtfyNotifiable interface {
	NtfyOptions() (NtfyOptions, error)
}

type NtfyOptions struct {
	Auth string `json:"auth,omitempty"`
}

func New(baseURL string, topic string) *NtfyChannel {
	return &NtfyChannel{
		baseURL: baseURL,
		topic:   topic,
		client:  &http.Client{},
	}
}

func (n *NtfyChannel) Notify(ctx context.Context, notifiable notify.Notifiable, notif notify.Notification) error {
	body, err := n.notificationToBody(notifiable, notif)
	if err != nil {
		return err
	}

	// set the topic if it's not set, default to the channel topic
	if body.Topic == "" {
		body.Topic = n.topic
	}

	if body.Topic == "" {
		return fmt.Errorf("a topic must be set for ntfy notifications")
	}

	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(body); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, n.baseURL, b)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	if notifiable, ok := notifiable.(NtfyNotifiable); ok {
		opts, err := notifiable.NtfyOptions()
		if err != nil {
			return err
		}

		if opts.Auth != "" {
			req.Header.Set("Authorization", opts.Auth)
		}
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from ntfy: %d", resp.StatusCode)
	}

	return nil
}

func (n *NtfyChannel) notificationToBody(notifiable notify.Notifiable, notif notify.Notification) (Message, error) {
	if ntfyNotif, ok := notif.(NtfyNotification); ok {
		return ntfyNotif.ToNtfy(notifiable)
	}

	return Message{
		Title:   notif.Subject(),
		Message: notif.Content(),
	}, nil
}
