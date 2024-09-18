# Resend channel

A simple channel for sending emails through [Resend](https://resend.com/)

## Interfaces

### Notification

| Function | Required | Description |
| --- | --- | --- |
| `ToResendMessage` | Yes | Converts the notification to a `resend.SendMessageRequest` struct. |

### Notifiable

| Function | Required | Description |
| --- | --- | --- |
| `ResendAddress` | Yes | Returns the email address to send the notification to. |

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
