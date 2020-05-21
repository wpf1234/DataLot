package models

type Message struct {
	MsgId    int     `json:"msg_id"`
	UserId   int     `json:"user_id"` // 进行操作的用户ID
	Username string  `json:"username"`
	Head     string  `json:"head"`
	Operate  string  `json:"operate"`
	Tm       int64   `json:"tm"`
	Dynamic  Dynamic `json:"dynamic"`
}

type MyComm struct {
	Id       int     `json:"id"`
	Username string  `json:"username"`
	Head     string  `json:"head"`
	Comment  string  `json:"comment"`
	Tm       int64   `json:"tm"`
	Dynamic  Dynamic `json:"dynamic"`
}
