package models

type Login struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Meid     string `json:"meid"`
	Desc     string `json:"desc"`
}

type Register struct {
	Username  string   `json:"username"`
	Phone     string   `json:"phone"`
	Password  string   `json:"password"`
	Interest  []string `json:"interest"`
}

type Reset struct {
	Phone       string `json:"phone"`
	NewPassword string `json:"new_password"`
}
