.PHONY: update_proto, install_protoc

run:
	go build
	./gameserver -port 4000

install_protoc:
	sudo aptitude install protobuf-compiler

update_proto:
	cd protobuf && protoc --go_out=. *.proto
