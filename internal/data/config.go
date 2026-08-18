package data

type ConfigData struct {
	Version string       `json:"version"`
	Server  ServerConfig `json:"server"`
}

type ServerConfig struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}
