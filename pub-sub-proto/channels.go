package protocol

const (
	// ChannelNameSubscribedChannels `/subscribed-channels` List of all currently subscribed channels
	ChannelNameSubscribedChannels = "/subscribed-channels"

	// ChannelNameStops `/stops` List of all stops
	ChannelNameStops = "/stops"

	// ChannelNameStop `/stop/<provider>/<stop-id>` Details of a specific stop
	ChannelNameStop = "/stop/%s/%s"

	// ChannelNameVehicles List of all vehicles
	ChannelNameVehicles = "/vehicles"

	// ChannelNameVehicle `/vehicle/<provider>/<vehicle-id>` Details of a specific vehicle
	ChannelNameVehicle = "/vehicle/%s/%s"

	// ChannelNameTrip `/trip/<provider>/<trip-id>` Details of a specific trip
	ChannelNameTrip = "/trip/%s/%s"

	// ChannelNameRoute `/route/<provider>/<route-id>` Details of a specific route
	ChannelNameRoute = "/route/%s/%s"
)
