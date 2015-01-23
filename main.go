package main

import (
	"flag"
	"log"
	"net/http"
)

var configFile = flag.String("config", "config.json", "configuration file")

func main() {
	flag.Parse()

	config, err := ReadConfiguration()
	if err != nil {
		log.Fatal("error starting: ", err)
	}

	Routes()

	if config.UseTLS {
		log.Println("Starting Secure WebSocket server at https://" + config.ServerAddress)
		log.Fatal(http.ListenAndServeTLS(config.ServerAddress, config.CertificateFile, config.KeyFile, nil))
	} else {
		log.Println("Starting WebSocket server at http://" + config.ServerAddress)
		log.Fatal(http.ListenAndServe(config.ServerAddress, nil))
	}
}
