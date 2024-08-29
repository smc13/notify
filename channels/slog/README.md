# Slog Channel
A simple channel for logging notifications to slog, predominantly used for debugging purposes.



# Usage
```go
import (
	"log/slog"
	"github.com/smc13/notify"
	slog_channel "github.com/smc13/notify/channels/slog"
)

service := notify.New(slog_channel.New(slog.Default()))

ctx := context.Background()
// All information in the context will be logged
ctx = context.WithValue(ctx, "key", "value")

// Send the notification through the slog channel
service.Notify(ctx, user, notification)
```
