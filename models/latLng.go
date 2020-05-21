package models

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
	Tm  int64 `json:"tm"`
}
