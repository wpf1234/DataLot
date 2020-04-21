package models

type User struct {
	Id        int      `json:"id"`
	Username  string   `json:"username"`
	Phone     string   `json:"phone"`
	Meid      string   `json:"meid"`
	PhoneDesc string   `json:"phone_desc"`
	Head      string   `json:"head"`
	Interest  []string `json:"interest"`
	UserSig   string   `json:"user_sig"`
}

// 高德地图相关信息
//type MapRes struct {
//	ServiceId  int `json:"service_id"`
//	TerminalId int `json:"terminal_id"`
//	TrackId    int `json:"track_id"`
//}

// 百度地图相关信息
type MapRes struct {
	ServiceId  int64  `json:"service_id"`
	EntityName string `json:"entity_name"`
}

type LoginRes struct {
	User  User   `json:"user"`
	Map   MapRes `json:"map"`
	Token string `json:"token"`
}
