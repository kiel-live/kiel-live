package protocol

const (
	// ___ Ctrl ___

	// List of all currently subscribed subjects of all clients
	SubjectSubscriptions = "ctrl.subscriptions"

	// Request to subscribe to a subject
	SubjectRequestSubscribe = "ctrl.subscribe"

	// Request to unsubscribe from a subject
	SubjectRequestUnsubscribe = "ctrl.unsubscribe"

	// Request to get the cache of a subject
	SubjectRequestCache = "ctrl.cache.request"

	// ___ Map ___

	// Single stops with gps position (data.map.<longitude>.<latitude>.stop.<id>)
	SubjectMapStop = "data.map.stop.%s"

	// Single vehicle with gps position (data.map.<longitude>.<latitude>.vehicle.<id>)
	SubjectMapVehicle = "data.map.vehicle.%s"

	// ___ Details ___

	// Details of a specific stop (data.stop.<id>)
	SubjectDetailsStop = "data.stop.%s"

	// Details of a specific vehicle (data.vehicle.<id>)
	SubjectDetailsVehicle = "data.vehicle.%s"

	// Details of a specific trip (data.trip.<id>)
	SubjectDetailsTrip = "data.map.trip.%s"

	// Details of a specific route (data.route.<id>)
	SubjectDetailsRoute = "data.route.%s"
)

const (
	// used to republish data that is not updated before being removed from the cache
	MaxCacheAge = 60 * 10 // 10 minutes

	DeletePayload = "---"
)
