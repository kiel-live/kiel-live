package protocol

const (
	// ___ Ctrl ___

	// List of all currently subscribed topics of all clients
	TopicSubscriptions = "ctrl.subscriptions"

	// Request to subscribe to a topic
	TopicRequestSubscribe = "ctrl.subscribe"

	// Request to unsubscribe from a
	TopicRequestUnsubscribe = "ctrl.unsubscribe"

	// Request to get the cache of a topic
	TopicRequestCache = "ctrl.cache.request"

	// ___ Map ___

	// Single stops with gps position (data.map.<longitude>.<latitude>.stop.<id>)
	TopicMapStop = "data.map.stop.%s"

	// Single vehicle with gps position (data.map.<longitude>.<latitude>.vehicle.<id>)
	TopicMapVehicle = "data.map.vehicle.%s"

	// ___ Details ___

	// Details of a specific stop (data.stop.<id>)
	TopicDetailsStop = "data.stop.%s"

	// Details of a specific vehicle (data.vehicle.<id>)
	TopicDetailsVehicle = "data.vehicle.%s"

	// Details of a specific trip (data.trip.<id>)
	TopicDetailsTrip = "data.map.trip.%s"

	// Details of a specific route (data.route.<id>)
	TopicDetailsRoute = "data.route.%s"
)

const (
	// used to republish data that is not updated before being removed from the cache
	MaxCacheAge = 60 * 10 // 10 minutes

	DeletePayload = "---"
)
