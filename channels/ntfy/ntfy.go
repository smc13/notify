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
	ToNtfy() (Message, error)
}

func New(baseURL string, topic string) *NtfyChannel {
	return &NtfyChannel{
		baseURL: baseURL,
		topic:   topic,
		client:  &http.Client{},
	}
}

func (n *NtfyChannel) Notify(ctx context.Context, auth []string, notif notify.Notification) error {
	body, err := n.notificationToBody(notif)
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
	if len(auth) > 0 {
		req.Header.Set("Authorization", auth[0])
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

func (n *NtfyChannel) notificationToBody(notif notify.Notification) (Message, error) {
	if ntfyNotif, ok := notif.(NtfyNotification); ok {
		return ntfyNotif.ToNtfy()
	}

	return Message{
		Title:   notif.Subject(),
		Message: notif.Content(),
	}, nil
}
