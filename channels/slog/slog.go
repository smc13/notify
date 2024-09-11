package slog

import (
	"context"
	"log/slog"

	"github.com/smc13/notify"
)

type SlogChannel struct {
	logger   *slog.Logger
	logLevel slog.Level
	logMsg   string
}

type SlogNotification interface {
	ToSlog(notify.Notifiable) []slog.Attr
}

func New(logger *slog.Logger) *SlogChannel {
	return &SlogChannel{
		logger:   logger,
		logLevel: slog.LevelInfo,
		logMsg:   "notification triggered",
	}
}

// SetLogLevel sets the log level for the notification
func (c *SlogChannel) SetLogLevel(level slog.Level) {
	c.logLevel = level
}

// SetLogMsg sets the message to be logged when a notification is triggered
func (c *SlogChannel) SetLogMsg(msg string) {
	c.logMsg = msg
}

// Notify logs the notification to the slog logger
func (c *SlogChannel) Notify(ctx context.Context, notifiable notify.Notifiable, notif notify.Notification) error {
	attrs := []slog.Attr{slog.String("kind", notif.Kind()), slog.Any("notifiable", notifiable)}
	if slogNotif, ok := notif.(SlogNotification); ok {
		attrs = append(attrs, slogNotif.ToSlog(notifiable)...)
	}

	c.logger.LogAttrs(ctx, c.logLevel, c.logMsg, attrs...)
	return nil
}
