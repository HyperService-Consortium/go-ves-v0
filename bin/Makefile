
client:
	go build -o ./binary/ves-client ./ves-client

server:
	go build -o ./binary/ves-server ./ves-server
	go build -o ./binary/central-ves-server ./central-ves-server
	if [ ! -f "./binary/ves-server-config.toml" ]; then cp ./ves-server/default-ves-server-config.toml ./binary/ves-server-config.toml; fi;

up:
	./binary/ves-server --config ./binary/ves-server-config.toml

up-central:
	./binary/central-ves-server

login:
	./binary/ves-client

.PHONY: client server up up-central
