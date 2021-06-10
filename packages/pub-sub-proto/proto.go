package protocol

// Message types used between server and client.
const (
	// Client: start authentication
	AuthMessage = "auth"

	// Server: Authentication succeeded
	AuthOKMessage = "auth-ok"

	// Server: Authentication failed
	AuthFailedMessage = "auth-error"

	// Client: Subscribe to channel
	SubscribeMessage = "subscribe"

	// Server: Subscribe succeeded
	SubscribeOKMessage = "subscribe-ok"

	// Server: Subscribe failed
	SubscribeErrorMessage = "subscribe-error"

	// Client: Publish message
	PublishMessage = "publish"

	// Publish: Publish succeeded
	PublishOKMessage = "publish-ok"

	// Publish: Publish failed
	PublishErrorMessage = "publish-error"

	// Server: Channel message
	ChannelMessage = "channel-message"

	// Client: Unsubscribe from channel
	UnsubscribeMessage = "unsubscribe"

	// Server: Unsubscribe succeeded
	UnsubscribeOKMessage = "unsubscribe-ok"

	// Server: Unsubscribe failed
	UnsubscribeErrorMessage = "unsubscribe-error"

	// Client: Send me more messages
	PollMessage = "poll"

	// Client: I'm still alive
	PingMessage = "ping"

	// Server: Unknown message
	UnknownMessage = "unknown"

	// Server: Server error
	ServerErrorMessage = "server-error"
)
