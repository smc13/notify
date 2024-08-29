package ntfy

import (
	"strconv"
	"time"
)

type Message struct {
	Topic    string          `json:"topic"`
	Message  string          `json:"message,omitempty"`
	Title    string          `json:"title,omitempty"`
	Tags     []string        `json:"tags,omitempty"`
	Priority MessagePriority `json:"priority,omitempty"`
	Actions  []MessageAction `json:"actions,omitempty"`
	Click    string          `json:"click,omitempty"`
	Attach   string          `json:"attach,omitempty"`
	Markdown bool            `json:"markdown,omitempty"`
	Icon     string          `json:"icon,omitempty"`
	Filename string          `json:"filename,omitempty"`
	Delay    MessageDelay    `json:"delay,omitempty"`
	Email    string          `json:"email,omitempty"`
	Call     string          `json:"call,omitempty"`
}

// MessageDelay is a time.Time that can be marshaled to JSON as a Unix timestamp
type MessageDelay time.Time

func (md MessageDelay) MarshalJSON() ([]byte, error) {
	return []byte(`"` + strconv.FormatInt(time.Time(md).Unix(), 10) + `"`), nil
}

type MessagePriority int

const (
	MessagePriorityLow    MessagePriority = 1
	MessagePriorityNormal MessagePriority = 3
	MessagePriorityHigh   MessagePriority = 5
)

type MessageActionType string

const (
	MessageActionTypeView      MessageActionType = "view"
	MessageActionTypeBroadcast MessageActionType = "broadcast"
	MessageActionTypeHttp      MessageActionType = "http"
)

type MessageAction struct {
	Action MessageActionType `json:"action"`
	Label  string            `json:"label"`
	Clear  bool              `json:"clear"`

	// Only used for view or http actions
	URL string `json:"url,omitempty"`
	// Only used for http actions
	Method string `json:"method,omitempty"`
	// Only used for http actions
	Headers map[string]string `json:"headers,omitempty"`
	// Only used for http actions
	Body string `json:"body,omitempty"`
	// Only used for broadcast actions
	Intent string `json:"intent,omitempty"`
	// ONly used for broadcast actions
	Extras map[string]string `json:"extras,omitempty"`
}
