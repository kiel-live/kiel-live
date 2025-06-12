package client

type HttpClient struct {
	url   string
	token string
}

func NewHttpClient(url, token string) Client {
	return &HttpClient{
		url:   url,
		token: token,
	}
}

func (h *HttpClient) Connect() error {
	return nil
}

func (h *HttpClient) Disconnect() error {
	return nil
}

func (h *HttpClient) IsConnected() bool {
	return false
}

func (h *HttpClient) Subscribe(topic string, cb SubscribeCallback) error {
	return nil
}

func (h *HttpClient) Unsubscribe(topic string) error {
	return nil
}

func (h *HttpClient) Publish(topic string, data string) error {
	return nil
}

func (h *HttpClient) SetConnectionHandler(handler func(connected bool)) {
}

func (h *HttpClient) SetTopicSubscriptionHandler(handler func(topics []string)) {

}
