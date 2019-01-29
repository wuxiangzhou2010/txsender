package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config is defalt config for tx_sender
type Config struct {
	Endpoints []string `json:"endpoints"`
	Rate      int      `json:"rate"`

	ChainAmount int  `json:"chainamount"`
	Silent      bool `json:"silent"`
	TxBuffer    int  `json:"txbuffer"`
	Worker      int  `json:"worker"`
}

// GetConfig get the config from json file
func GetConfig(path string) *Config {
	var config Config
	configFile, err := os.Open(path)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config
}
