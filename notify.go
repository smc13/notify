package notify

import (
	"context"
	"fmt"
	"sync"
)

type Notify struct {
	channels []Channel
}

func NewNotify(channels ...Channel) *Notify {
	return &Notify{channels: channels}
}

func (n *Notify) AddChannel(channel Channel) {
	n.channels = append(n.channels, channel)
}

// Notify sends a notification to all channels synchronously
// It will return the first error that occurred during the notification process
func (n *Notify) Notify(ctx context.Context, notifiable Notifiable, notif Notification) error {
	for _, channel := range n.channels {
		if err := n.sendToChannel(ctx, channel, notifiable, notif); err != nil {
			return err
		}
	}

	return nil
}

type NotifyResult struct {
	Errors []error
	wg     sync.WaitGroup
}

func (r *NotifyResult) Wait() {
	r.wg.Wait()
}

func (r *NotifyResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *NotifyResult) Error() error {
	if r.HasErrors() {
		return fmt.Errorf("notify: %d errors occurred", len(r.Errors))
	}

	return nil
}

// NotifyConcurrent sends a notification to all channels asynchronously
// It returns a slice of errors that occurred during the notification process, the order of the errors is not guaranteed
func (n *Notify) NotifyConcurrent(ctx context.Context, notifiable Notifiable, notif Notification) *NotifyResult {
	result := &NotifyResult{
		wg:     sync.WaitGroup{},
		Errors: make([]error, 0, len(n.channels)),
	}

	result.wg.Add(len(n.channels))

	for _, channel := range n.channels {
		go func(channel Channel) {
			defer result.wg.Done()
			if err := n.sendToChannel(ctx, channel, notifiable, notif); err != nil {
				result.Errors = append(result.Errors, err)
			}
		}(channel)
	}

	return result
}

func (n *Notify) sendToChannel(ctx context.Context, channel Channel, notifiable Notifiable, notif Notification) error {
	shouldSend := notifiable.ShouldSendNotification(ctx, channel, notif)
	if !shouldSend {
		return nil
	}

	return channel.Notify(ctx, notifiable, notif)
}
