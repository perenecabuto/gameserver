package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() {
	router := mux.NewRouter()

	router.PathPrefix("/protobuf/").Handler(http.StripPrefix("/protobuf/", http.FileServer(http.Dir("protobuf/"))))

	router.Handle("/ws/chat", NewChatServer())
	router.Handle("/ws/game", NewGameServer())
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("webroot/")))

	http.Handle("/", router)
}
