package models

type Report struct {
	Id     int      `json:"id"`
	Reason []string `json:"reason"`
}
