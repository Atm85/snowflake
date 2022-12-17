package models

type Snowdepth []struct {
	StationID         string `json:"station_id"`
	StationName       string `json:"station_name"`
	ProviderID        string `json:"provider_id"`
	ProviderStationID string `json:"provider_station_id"`
	MaxSnowDepth      string `json:"max_snow_depth_in"`
}
