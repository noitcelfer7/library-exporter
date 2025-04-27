package config

type Config struct {
	Grpc struct {
		Server struct {
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"server"`
	} `json:"grpc"`
	Http struct {
		Server struct {
			Host string `json:"host"`
			Port string `json:"port"`
		} `json:"server"`
	} `json:"http"`
	Postgresql struct {
		Url string `json:"url"`
	} `json:"postgresql"`
}
