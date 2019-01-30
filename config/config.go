package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config is default config for tx_sender
type Config struct {
	Endpoints []string `json:"endpoints"`
	Rate      int      `json:"rate"`

	ChainAmount int  `json:"chainAmount"`
	Silent      bool `json:"silent"`
	TxBuffer    int  `json:"txBuffer"`
	Worker      int  `json:"worker"`
}

// GetConfig get the config from json file
func getConfig(path string) *Config {
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

func GetConfig() *Config {

	cfg := getConfig("./config.json")
	PrintConfig(cfg)
	return cfg
}

func PrintConfig(cfg *Config) {
	log.Print(
		"\n#######################\n",
		"rate\t\t\t: ", cfg.Rate,
		"\nendpoint\t\t: ", cfg.Endpoints,
		"\ntxBuffer\t\t: ", cfg.TxBuffer,
		"\nworker\t\t\t: ", cfg.Worker,
		"\n#######################\n",
	)
}
