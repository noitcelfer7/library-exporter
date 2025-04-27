package config

type Config struct {
	Grpc struct {
		Server struct {
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"server"`
	} `json:"grpc"`
}
