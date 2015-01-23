package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	ServerAddress   string `json:"serverAddress"`
	CertificateFile string `json:"certFile"`
	KeyFile         string `json:"keyFile"`
	UseTLS          bool   `json:"useTLS"`
}

func ReadConfiguration() (conf *Configuration, err error) {
	var file *os.File
	if file, err = os.Open(*configFile); err == nil {
		log.Println("Reading config:", *configFile)

		conf = &Configuration{}
		if err = json.NewDecoder(file).Decode(&conf); err != nil {
			return
		}
	}

	return
}
