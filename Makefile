.PHONY: update_proto, install_protoc

run:
	@go build
	@./gameserver

install_protoc_linux:
	@sudo aptitude install protobuf-compiler

install_protoc:
	@get -u github.com/golang/protobuf/{proto,protoc-gen-go}

update_proto:
	@cd protobuf && protoc --go_out=. *.proto
