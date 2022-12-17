package models

type Snowfall []struct {
	StationID         string `json:"station_id"`
	StationName       string `json:"station_name"`
	ProviderID        string `json:"provider_id"`
	ProviderStationID string `json:"provider_station_id"`
	MaxSnowFall       string `json:"max_snowfall_24h_in"`
}
