build: test
	go build -o ./bin/adminclient mathbattle/cmd/adminclient
	go build -o ./bin/tgbot mathbattle/cmd/tgbot
	go build -o ./bin/sendmsg mathbattle/cmd/sendmsg
test:
	go test ./...
