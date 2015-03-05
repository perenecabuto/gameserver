package main

import (
	"os"
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

	address := config.ServerAddress
	if len(os.Getenv("PORT")) > 0 {
        address = "0.0.0.0:" + os.Getenv("PORT")
	}

	Routes()

	if config.UseTLS {
		log.Println("Starting Secure WebSocket server at https://" + config.ServerAddress)
		log.Fatal(http.ListenAndServeTLS(config.ServerAddress, config.CertificateFile, config.KeyFile, nil))
	} else {
		log.Println("Starting WebSocket server at http://" + address)
		log.Fatal(http.ListenAndServe(address, nil))
	}
}
