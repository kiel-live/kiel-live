package protocol

type Location struct {
	Longitude float32 `json:"longitude"` // exp: 54.306
	Latitude  float32 `json:"latitude"`  // exp: 10.149
	Heading   int     `json:"heading"`   // in degree
}
