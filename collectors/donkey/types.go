package main

type DonkeyResponse struct {
	Hubs []struct {
		HubType                string `json:"hub_type"`
		ID                     string `json:"id"`
		AccountID              int    `json:"account_id"`
		Name                   string `json:"name"`
		Latitude               string `json:"latitude"`
		Longitude              string `json:"longitude"`
		Radius                 int    `json:"radius"`
		MaximumCapacity        int    `json:"maximum_capacity"`
		AvailableVehiclesCount int    `json:"available_vehicles_count"`
		VehiclesCount          int    `json:"vehicles_count"`
		AvailableVehicles      struct {
			Bike     int `json:"bike"`
			Ebike    int `json:"ebike"`
			Escooter int `json:"escooter"`
			Trailer  int `json:"trailer"`
			Cargo    int `json:"cargo"`
			Ecargo   int `json:"ecargo"`
		} `json:"available_vehicles"`
	} `json:"hubs"`
	Accounts []struct {
		ID                    int           `json:"id"`
		SupportedVehicleTypes []string      `json:"supported_vehicle_types"`
		CollectConsent        bool          `json:"collect_consent"`
		DayDeals              []interface{} `json:"day_deals"`
		Currency              string        `json:"currency"`
		RelocationFees        struct {
			Num500   string `json:"500"`
			Num5000  string `json:"5000"`
			Num50000 string `json:"50000"`
		} `json:"relocation_fees"`
		LostVehicleFees struct {
			Bike     string `json:"bike"`
			Ebike    string `json:"ebike"`
			Escooter string `json:"escooter"`
			Trailer  string `json:"trailer"`
			Cargo    string `json:"cargo"`
			Ecargo   string `json:"ecargo"`
		} `json:"lost_vehicle_fees"`
		LostInsuredVehicleFees struct {
			Bike   string `json:"bike"`
			Ebike  string `json:"ebike"`
			Cargo  string `json:"cargo"`
			Ecargo string `json:"ecargo"`
		} `json:"lost_insured_vehicle_fees"`
		OtherFees struct {
			BikeLost1Months string `json:"bike_lost_1_months"`
			BikeLost2Months string `json:"bike_lost_2_months"`
			DebtInterest    string `json:"debt_interest"`
		} `json:"other_fees"`
	} `json:"accounts"`
	Schedules []interface{} `json:"schedules"`
}
