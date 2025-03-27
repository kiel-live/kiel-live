package protocol

type Action struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type ActionType string

const (
	ActionTypeNavigateTo ActionType = "navigate-to"
	ActionTypeRent       ActionType = "rent"
)
