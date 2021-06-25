package api

import protocol "github.com/kiel-live/kiel-live/packages/pub-sub-proto"

type route struct {
	Stops         []tripStops `json:"actual"`
	OldStops      []tripStops `json:"old"`
	DirectionText string      `json:"directionText"`
	RouteName     string      `json:"routeName"`
}

func (r *route) parse() *protocol.Route {
	// TODO
	return nil
}

func GetRoute(routeID string) route {
	// TODO
	return route{}
}
