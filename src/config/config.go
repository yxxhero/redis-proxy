package config

import (
	"encoding/json"
	"log"
	"os"
	"structure"
)

type Config struct {
	Service      string              `json:"service"`
	Host         string              `json:"host"`
	Port         uint16              `json:"port"`
	WebPort      uint16              `json:"webport"`
	Strategy     string              `json:"strategy"`
	Heartbeat    int                 `json:"heartbeat"`
	MaxCointegration int                 `json:"maxcointegration"`
	MaxConnection int                 `json:"maxconnection"`
	MasterHost   string              `json:"masterhost"`
	MasterPort   uint16              `json:"masterport"`
	Backends     []structure.Backend `json:"backends"`
}

func Load(filename string) (*Config, error) {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		log.Println("load config failed:", err)
	} else {
		buff := make([]byte, 1024)
		end, _ := file.Read(buff)
		err = json.Unmarshal(buff[:end], &config)
		if err != nil {
			log.Println("decode json config failed:", err)
		}
	}
	log.Println("success load config file:", filename)
	return &config, err
}
