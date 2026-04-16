package types

type StationStatus struct {
	StationID             string             `json:"station_id"`
	NumVehiclesAvailable  int                `json:"num_vehicles_available"`
	NumDocksAvailable     int                `json:"num_docks_available"`
	IsInstalled           bool               `json:"is_installed"`
	IsRenting             bool               `json:"is_renting"`
	IsReturning           bool               `json:"is_returning"`
	LastReported          string             `json:"last_reported"`
	VehicleTypesAvailable []VehicleTypeCount `json:"vehicle_types_available,omitempty"`
}

type VehicleTypeCount struct {
	VehicleTypeID string `json:"vehicle_type_id"`
	Count         int    `json:"count"`
}

type StationInformation struct {
	StationID        string            `json:"station_id"`
	Name             []LocalizedString `json:"name"`
	Lat              float64           `json:"lat"`
	Lon              float64           `json:"lon"`
	RegionID         string            `json:"region_id,omitempty"`
	IsVirtualStation bool              `json:"is_virtual_station,omitempty"`
	Capacity         int               `json:"capacity,omitempty"`
	RentalURIs       RentalURIs        `json:"rental_uris"`
}

type LocalizedString struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}

type RentalURIs struct {
	Android string `json:"android,omitempty"`
	IOS     string `json:"ios,omitempty"`
}

type Response struct {
	TTL         int    `json:"ttl"`
	LastUpdated string `json:"last_updated"`
	Data        any    `json:"data"`
}

type VehicleStatus struct {
	VehicleID          string     `json:"vehicle_id"`
	Lat                float64    `json:"lat"`
	Lon                float64    `json:"lon"`
	IsReserved         bool       `json:"is_reserved"`
	IsDisabled         bool       `json:"is_disabled"`
	LastReported       string     `json:"last_reported"`
	RentalURIs         RentalURIs `json:"rental_uris"`
	VehicleTypeID      string     `json:"vehicle_type_id,omitempty"`
	CurrentRangeMeters int        `json:"current_range_meters,omitempty"`
	CurrentFuelPercent float64    `json:"current_fuel_percent,omitempty"`
	StationID          string     `json:"station_id,omitempty"`
}
