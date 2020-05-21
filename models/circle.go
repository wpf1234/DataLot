package models

type Content struct {
	Text    string   `json:"text"`
	Picture []string `json:"picture"`
	Auth    int      `json:"auth"`
}

type Comment struct {
	DynamicId int    `json:"dynamic_id"`
	Id        int    `json:"id"` //用户ID
	CommUser  string `json:"comm_user"`
	Head      string `json:"head"`
	Context   string `json:"context"`
	Tm        int64  `json:"tm"`
}

type WriteComm struct {
	DynamicId int    `json:"dynamic_id"`
	UserId    int    `json:"user_id"`
	CommId    int    `json:"comm_id"`
	Context   string `json:"context"`
}

type Reply struct {
	Id        int    `json:"id"` // 用户ID
	ReplyUser string `json:"reply_user"`
	Head      string `json:"head"`
	Context   string `json:"context"`
	Tm        int64  `json:"tm"`
}

type WriteReply struct {
	CommId  int    `json:"comm_id"`
	UserId  int    `json:"user_id"`
	Context string `json:"context"`
	Tm      int64  `json:"tm"`
}

type Dynamic struct {
	Id       int       `json:"id"`
	UserId   int       `json:"user_id"`
	Username string    `json:"username"`
	Head     string    `json:"head"`
	Content  string    `json:"content"`
	Picture  []string  `json:"picture"`
	Tm       int64     `json:"tm"`
	Favorite int       `json:"favorite"`
	Comment  int       `json:"comment"`
	CommList []Comment `json:"comm_list"`
	Like     bool      `json:"like"`
}

type MyDynamic struct {
	Id       int       `json:"id"`
	Username string    `json:"username"`
	Head     string    `json:"head"`
	Content  string    `json:"content"`
	Picture  []string  `json:"picture"`
	Tm       int64     `json:"tm"`
	Favorite int       `json:"favorite"`
	Comment  int       `json:"comment"`
	CommList []Comment `json:"comm_list"`
}

type Circle struct {
	Recommend []Dynamic `json:"recommend"`
	Friend    []Dynamic `json:"friend"`
	News      []Dynamic `json:"news"`
}

type UploadPIC struct {
	UserId  int    `json:"user_id"`
	Picture string `json:"picture"`
}
