# Notify

Notify is a simple notification service that allows you to send notifications to multiple channels.
It's designed to be dead simple to use and easy to extend.

## Installation
```bash
go get github.com/smc13/notify
```

## Getting Started

### Set up the service
```go
import (
	"github.com/smc13/notify"
	"github.com/smc13/notify/channels/discord"
)

discordChannel := discord.NewWebhookChannel(discord.WebhookURLToParts("https://discord.com/api/webhooks/1234567890/ABCDEFGHIJKLMN"))

service := notify.New(discordChannel)
```

### Create a Notifiable
A `Notifiable` is a struct with a `RouteNotificationFor` function that can be used to
customize channel notifications for the struct.
It could be used to mention a user by their Discord ID or skip sending to a channel based on user preferences.

#### Grouped Notifiable
A `GroupedNotifiable` is provided to allow sending a single notification to multiple notifiables.
How this is handled is up to the channel implementation; e.g. Discord webhooks will send a single message with multiple mentions
while Telegram will notify each chat individually.

```go
type User struct {
	ID 								int64
	Email 						string
	DiscordID 				string
	EnableTelegram 		bool
}

func (u User) RouteNotificationFor(ctx context.Context, channel notify.Channel, notif notify.Notification) ([]string, error) {
	switch channel.(type) {
		// Mention the discord user by their ID
		case discord.WebhookChannel:
			return []string{u.DiscordID}, nil
		// Skip sending to telegram if the user has it disabled
		case telegram.TelegramChannel:
			if !u.EnableTelegram {
				return nil, notify.ErrSkipNotification
			}
	}

	// No special routing for this channel but still send
	return nil, nil
}
```

### Create a Notification
A `Notification` is a struct that contains the information to be sent to the channels.
At a minimum it requires:
1. A `Kind` function that identifies the type of notification.
2. A `Subject` function that provides a short description of the notification.
3. A `Content` function that provides the full content of the notification.

> [!NOTE]
> Some channels may allow for more customization of the notification through their own interfaces.
> Be sure to check the documentation for the channels you are using.

```go
type UserRegisteredNotification struct {
	User User
}

func (n UserRegisteredNotification) Kind() string { return "user_registered" }
func (n UserRegisteredNotification) Subject() string { return "User Registered" }
func (n UserRegisteredNotification) Content() string {
	return fmt.Sprintf("User %s has registered", n.User.Email)
}
```

### Send a Notification
```go
ctx := context.Background()

newUser := User{ID: 1, Email: "test@email.com"}
adminUser := User{ID: 2, DiscordID: "1234567890"}

// Send the user registered notification to the admin user through the discord channel we set up earlier
service.Notify(ctx, adminUser, UserRegisteredNotification{User: newUser})
```

## Custom Channels
Custom channels can be created by implementing the `Channel` interface.

```go
type CustomChannel struct {
	notfications []struct { notification notify.Notification, notifiable notify.Notifiable }
}

func (c CustomChannel) Notify(ctx context.Context, notifiable notify.Notifiable, notification notify.Notification) error {
	// for example, store the notification to be sent later
	c.notifications = append(c.notifications, struct { notification notify.Notification, notifiable notify.Notifiable }{notification, notifiable})

	return nil
}
```


## Todo
- [ ] Additional channels
	- [ ] Gotify
	- [ ] Slack
 	- [ ] Resend
  - [ ] SMTP
  - [ ] Telgram
  - [ ] Vonage
  - [ ] Webhook
- [ ] Add batch notification support
