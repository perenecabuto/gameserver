package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Routes() *mux.Router {
	router := mux.NewRouter()

	router.PathPrefix("/protobuf/").Handler(http.StripPrefix("/protobuf/", http.FileServer(http.Dir("protobuf/"))))

	router.Handle("/ws/chat", NewChatServer())
	router.Handle("/ws/game", GameHandler())
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("webroot/")))

	return router
}
