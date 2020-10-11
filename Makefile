build: test
	go build -o ./bin/adminclient mathbattle/cmd/adminclient
	go build -o ./bin/tgbot mathbattle/cmd/tgbot
test:
	go test ./...