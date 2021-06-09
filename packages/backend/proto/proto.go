package proto

// Message types used between server and client.
const (
	// Client: start authentication
	AuthMessage = "auth"

	// Server: Authentication succeeded
	AuthOKMessage = "authOk"

	// Server: Authentication failed
	AuthFailedMessage = "authError"

	// Client: Subscribe to channel
	SubscribeMessage = "subscribe"

	// Server: Subscribe succeeded
	SubscribeOKMessage = "subscribeOk"

	// Server: Subscribe failed
	SubscribeErrorMessage = "subscribeError"

	// Server: Broadcast message
	MessageMessage = "message"

	// Client: Unsubscribe from channel
	UnsubscribeMessage = "unsubscribe"

	// Server: Unsubscribe succeeded
	UnsubscribeOKMessage = "unsubscribeOk"

	// Server: Unsubscribe failed
	UnsubscribeErrorMessage = "unsubscribeError"

	// Client: Send me more messages
	PollMessage = "poll"

	// Client: I'm still alive
	PingMessage = "ping"

	// Server: Unknown message
	UnknownMessage = "unknown"

	// Server: Server error
	ServerErrorMessage = "serverError"
)
