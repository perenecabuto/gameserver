package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Configuration struct {
	Address         string `json:"addr"`
	Port            string `json:"port"`
	CertificateFile string `json:"certFile"`
	KeyFile         string `json:"keyFile"`
	UseTLS          bool   `json:"useTLS"`
}

var configFile = flag.String("config", "config.json", "configuration file")

func main() {
	flag.Parse()

	config, err := ReadConfiguration()
	if err != nil {
		log.Fatal("error starting: ", err)
	}

	serverAddress := fmt.Sprintf("%s:%s", config.Address, config.Port)

	Routes()

	if config.UseTLS {
		log.Println(fmt.Sprintf("Starting Secure WebSocket server at https://%s", serverAddress))
		log.Fatal(http.ListenAndServeTLS(serverAddress, config.CertificateFile, config.KeyFile, nil))
	} else {
		log.Println(fmt.Sprintf("Starting WebSocket server at http://%s", serverAddress))
		log.Fatal(http.ListenAndServe(serverAddress, nil))
	}
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

func Routes() {
	router := mux.NewRouter()

	router.PathPrefix("/protobuf/").Handler(http.StripPrefix("/protobuf/", http.FileServer(http.Dir("protobuf/"))))

	router.Handle("/ws/chat", NewChatServer())
	router.Handle("/ws/game", NewGameServer())
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("webroot/")))

	http.Handle("/", router)
}
