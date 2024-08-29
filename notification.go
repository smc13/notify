package notify

type Notification interface {
	// The kind of notification, should be unique. eg: "user.created", "invoice.paid", "order.shipped"
	Kind() string
	// Subject of the notification eg: "User Created", "Invoice Paid", "Order Shipped"
	Subject() string
	// Content of the notification
	Content() string
}
