package data

type ConfigData struct {
	Version string       `json:"version"`
	Server  ServerConfig `json:"server"`
	DB      DbConfig     `json:"db"`
	Auth    AuthConfig   `json:"auth"`
	Context string       `json:"context"`
}

type ServerConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type DbConfig struct {
	Machine  string `json:"machine"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type AuthConfig struct {
	URL   string `json:"url"`
	Right string `json:"right"`
}
