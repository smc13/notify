# SMTP channel

A simple channel for sending notifications through SMTP, using the [wneessen/go-mail](https://github.com/wneessen/go-mail) package.

## Interfaces

### Notification

| Function | Required | Description |
| --- | --- | --- |
| `ToSMTPMessage` | Yes | Converts the notification to a `mail.Msg` struct. |

### Notifiable

| Function | Required | Description |
| --- | --- | --- |
| `SMTPAddress` | Yes | Returns the email address to send the notification to. |

## Usage

```go
import (
	"github.com/smc13/notify"
	"github.com/smc13/notify/channels/smtp"
	"github.com/wneessen/go-mail"
)

c, err := mail.NewClient("smtp.example.com", mail.WithPort(25), mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithUsername("my_username"), mail.WithPassword("extremely_secret_pass"))
if err != nil {
	panic(err)
}

service := notify.NewNotify(smtp.New(c))
```
