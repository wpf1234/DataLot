package models

type UserInfo struct {
	Id       int       `json:"id"`
	Username string    `json:"username"`
	Head     string    `json:"head"`
	Dynamic  []Dynamic `json:"dynamic"`
}
