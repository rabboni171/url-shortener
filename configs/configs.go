package configs

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	AppParams     Params        `json:"app_params"`
	DBRedisParams DBRedisParams `json:"db_redis_params"`
}

type Params struct {
	ServerURL    string `json:"server_url"`
	PortRun      string `json:"port_run"`
	WriteTimeout int64  `json:"write_timeout"`
	ReadTimeout  int64  `json:"read_timeout"`
	BaseURL      string `json:"base_url"`
}

type DBRedisParams struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func InitConfig() (*Config, error) {
	configFile, err := os.Open("./configs.json")
	if err != nil {
		return nil, fmt.Errorf("couldn't open config file: %v", err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
		}
	}(configFile)

	var config Config
	if err = json.NewDecoder(configFile).Decode(&config); err != nil {
		return nil, fmt.Errorf("couldn't decode settings json file: %v", err)
	}

	return &config, nil
}
