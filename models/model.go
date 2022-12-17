package models

type Result struct {
	Feed      string `json:"feed"`
	Version   string `json:"version"`
	Snowdepth `json:"snow_depth"`
	Snowfall  `json:"snowfall_24h"`
}
