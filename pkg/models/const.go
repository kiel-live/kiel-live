package models

const (
	// List of all currently subscribed topics of all clients
	TopicSubscriptions string = "ctrl.subscriptions"

	// Request to subscribe to a topic
	TopicRequestSubscribe string = "ctrl.subscribe"

	// Request to unsubscribe from a
	TopicRequestUnsubscribe string = "ctrl.unsubscribe"

	// Request to get the cache of a topic
	TopicRequestCache string = "ctrl.cache.request"

	// Single stop (data.map.stop.<id>)
	TopicStop string = "data.map.stop.%s"

	// Single vehicle (data.map.vehicle.<id>)
	TopicVehicle string = "data.map.vehicle.%s"

	// Single trip (data.map.trip.<id>)
	TopicTrip string = "data.trip.%s"

	// Single route (data.route.<id>)
	TopicRoute string = "data.route.%s"
)

const (
	// used to republish data that is not updated before being removed from the cache
	MaxCacheAge = 60 * 10 // 10 minutes

	DeletePayload = "---"
)
