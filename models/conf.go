package models

type MysqlConf struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DB       string `json:"db"`
}

type LogConf struct {
	LogPath string `json:"log_path"`
	LogFile string `json:"log_file"`
}

// 高德地图配置
//type MapConf struct {
//	Key         string `json:"key"`
//	ServiceUrl  string `json:"service_url"`
//	TerminalUrl string `json:"terminal_url"`
//	TrackUrl    string `json:"track_url"`
//	SearchUrl   string `json:"search_url"`
//}

// 百度地图配置
type MapConf struct {
	ServiceId  int64  `json:"service_id"`
	Key        string `json:"key"`
	EntityUrl  string `json:"entity_url"`
	EntityList string `json:"entity_list"`
	AddPoint   string `json:"add_point"`
	TrackUrl   string `json:"track_url"`
}

// 腾讯 IM 即时通讯
type IMConf struct {
	Appid int64  `json:"appid"`
	Key   string `json:"key"`
}
