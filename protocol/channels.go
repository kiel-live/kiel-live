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
	SubjectMapStop = "data.map.%d.%d.stop.%s"

	// Single vehicle with gps position (data.map.<longitude>.<latitude>.vehicle.<id>)
	SubjectMapVehicle = "data.map.%d.%d.vehicle.%s"

	// ___ Details ___

	// Details of a specific stop (data.stop.<id>)
	SubjectDetailsStop = "data.stop.%s"

	// Details of a specific vehicle (data.vehicle.<id>)
	SubjectDetailsVehicle = "data.vehicle.%s"

	// Details of a specific trip (data.trip.<id>)
	SubjectDetailsTrip = "data.trip.%s"

	// Details of a specific route (data.route.<id>)
	SubjectDetailsRoute = "data.route.%s"
)
