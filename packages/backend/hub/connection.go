package hub

type connection interface {
	Send(channel, message string)
	Process(t string, args []string)
	GetToken() string
}
