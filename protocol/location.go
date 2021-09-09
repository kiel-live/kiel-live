package protocol

type Location struct {
	Longitude int `json:"longitude"` // exp: 54.306 * 3600000 = longitude
	Latitude  int `json:"latitude"`  // exp: 10.149 * 3600000 = latitude
	Heading   int `json:"heading"`   // in degree
}
