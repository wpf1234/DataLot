package models

type MapToken struct {
	Note        string   `json:"note"`
	Scopes      []string `json:"scopes"`
	AllowedUrls []string `json:"allowedUrls"`
}
