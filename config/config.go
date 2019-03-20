package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Config is default config for txsender
type Config struct {
	Endpoint string `json:"endpoint"`
	Rate     int32  `json:"rate"`

	ChainAmount    int   `json:"chainAmount"`
	Silent         bool  `json:"silent"`
	SignedTxBuffer int   `json:"signedTxBuffer"`
	RawTxBuffer    int   `json:"rawTxBuffer"`
	Last           int32 `json:"last"`
	TxPerRecipient int   `json:"txPerRecipient"`
	SignerWorker   int   `json:"signerWorker"`
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
		"\n##### config ########\n",
		"rate\t\t\t: ", cfg.Rate,
		"\nendpoint\t\t: ", cfg.Endpoint,
		"\ntxBuffer\t\t: ", cfg.RawTxBuffer,
		"\nsignedTxBuffer\t\t: ", cfg.SignedTxBuffer,
		"\nlast\t\t\t: ", cfg.Last,
		"\n#######################\n",
	)
}
