package hub

type connection interface {
	Send(channel string, message string)
	GetToken() string
}
