build:
	npm --prefix ./client install ./client

include .env
export $(shell sed 's/=.*//' .env)

run:
	npm --prefix ./client start > client-output.log &
	go run ./server/index.go > server-output.log &

run-server:
	go run ./server/index.go

run-console:
	go run ./clients/console/index.go

run-web:
	npm --prefix ./clients/web run build && npm --prefix ./clients/web start

install-web-dependencies:
	npm --prefix ./clients/web install

stop:
	kill $(lsof -t -i:$(SERVER_PORT)) 
	kill $(lsof -t -i:$(CLIENT_PORT)) 
