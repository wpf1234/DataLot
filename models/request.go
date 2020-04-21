package models

type Login struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type Register struct {
	Username  string   `json:"username"`
	Phone     string   `json:"phone"`
	Password  string   `json:"password"`
	Meid      string   `json:"meid"`
	PhoneDesc string   `json:"phone_desc"`
	Interest  []string `json:"interest"`
}

type Reset struct {
	Phone       string `json:"phone"`
	NewPassword string `json:"new_password"`
}
