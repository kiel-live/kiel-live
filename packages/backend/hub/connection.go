package hub

type connection interface {
	Send(channel, message string)
	GetToken() string
}
