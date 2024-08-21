package main

type NextbikeResponse struct {
	Countries []struct {
		Lat                   float64 `json:"lat"`
		Lng                   float64 `json:"lng"`
		Zoom                  int     `json:"zoom"`
		Name                  string  `json:"name"`
		Hotline               string  `json:"hotline"`
		Domain                string  `json:"domain"`
		Language              string  `json:"language"`
		Email                 string  `json:"email"`
		Timezone              string  `json:"timezone"`
		Currency              string  `json:"currency"`
		CountryCallingCode    string  `json:"country_calling_code"`
		SystemOperatorAddress string  `json:"system_operator_address"`
		Country               string  `json:"country"`
		CountryName           string  `json:"country_name"`
		Terms                 string  `json:"terms"`
		Policy                string  `json:"policy"`
		Website               string  `json:"website"`
		ShowBikeTypes         bool    `json:"show_bike_types"`
		ShowBikeTypeGroups    bool    `json:"show_bike_type_groups"`
		ShowFreeRacks         bool    `json:"show_free_racks"`
		BookedBikes           int     `json:"booked_bikes"`
		SetPointBikes         int     `json:"set_point_bikes"`
		AvailableBikes        int     `json:"available_bikes"`
		CappedAvailableBikes  bool    `json:"capped_available_bikes"`
		NoRegistration        bool    `json:"no_registration"`
		Pricing               string  `json:"pricing"`
		FaqURL                string  `json:"faq_url"`
		StoreURIAndroid       string  `json:"store_uri_android"`
		StoreURIIos           string  `json:"store_uri_ios"`
		Cities                []struct {
			UID         int     `json:"uid"`
			Lat         float64 `json:"lat"`
			Lng         float64 `json:"lng"`
			Zoom        int     `json:"zoom"`
			MapsIcon    string  `json:"maps_icon"`
			Alias       string  `json:"alias"`
			Break       bool    `json:"break"`
			Name        string  `json:"name"`
			NumPlaces   int     `json:"num_places"`
			RefreshRate string  `json:"refresh_rate"`
			Bounds      struct {
				SouthWest struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"south_west"`
				NorthEast struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"north_east"`
			} `json:"bounds"`
			BookedBikes          int  `json:"booked_bikes"`
			SetPointBikes        int  `json:"set_point_bikes"`
			AvailableBikes       int  `json:"available_bikes"`
			ReturnToOfficialOnly bool `json:"return_to_official_only"`
			BikeTypes            struct {
				Num28  int `json:"28"`
				Num29  int `json:"29"`
				Num71  int `json:"71"`
				Num150 int `json:"150"`
			} `json:"bike_types"`
			Website string `json:"website"`
			Places  []struct {
				UID                  int         `json:"uid"`
				Lat                  float64     `json:"lat"`
				Lng                  float64     `json:"lng"`
				Bike                 bool        `json:"bike"`
				Name                 string      `json:"name"`
				Address              interface{} `json:"address"`
				Spot                 bool        `json:"spot"`
				Number               int         `json:"number"`
				BookedBikes          int         `json:"booked_bikes"`
				Bikes                int         `json:"bikes"`
				BikesAvailableToRent int         `json:"bikes_available_to_rent"`
				BikeRacks            int         `json:"bike_racks"`
				FreeRacks            int         `json:"free_racks"`
				SpecialRacks         int         `json:"special_racks"`
				FreeSpecialRacks     int         `json:"free_special_racks"`
				Maintenance          bool        `json:"maintenance"`
				TerminalType         string      `json:"terminal_type"`
				BikeList             []struct {
					Number         string      `json:"number"`
					BikeType       int         `json:"bike_type"`
					LockTypes      []string    `json:"lock_types"`
					Active         bool        `json:"active"`
					State          string      `json:"state"`
					ElectricLock   bool        `json:"electric_lock"`
					Boardcomputer  int64       `json:"boardcomputer"`
					PedelecBattery interface{} `json:"pedelec_battery"`
					BatteryPack    interface{} `json:"battery_pack"`
				} `json:"bike_list"`
				BikeNumbers []string `json:"bike_numbers"`
				BikeTypes   struct {
				} `json:"bike_types,omitempty"`
				PlaceType string `json:"place_type"`
				RackLocks bool   `json:"rack_locks"`
			} `json:"places"`
		} `json:"cities"`
	} `json:"countries"`
}
