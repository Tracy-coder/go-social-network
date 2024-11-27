idls:
	hz update -idl ./idl/network.proto

run:
	go build && ./go-social-network