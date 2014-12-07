package main

import (
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/perenecabuto/gameserver/protobuf"
	"log"
	"net/http"
)

var (
	port = flag.String("port", "8888", "Define what TCP port to bind to")
	root = flag.String("root", ".", "Define the root filesystem path")
)

func main() {
	flag.Parse()

	object := &protobuf.Object{
		Id:       proto.Int32(1),
		Position: &protobuf.Object_Position{X: proto.Int32(1), Y: proto.Int32(1)},
		Status:   protobuf.Object_IDLE.Enum(),
	}

	data, err := proto.Marshal(object)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	fmt.Println(len(data))

	panic(http.ListenAndServe(":"+*port, http.FileServer(http.Dir(*root))))
}
