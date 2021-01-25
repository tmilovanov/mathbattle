build:
	go build -o ./bin/mb-server mathbattle/cmd/mb-server
	go build -o ./bin/mb-bot mathbattle/cmd/mb-bot
	go build -o ./bin/mb-admin mathbattle/cmd/mb-admin
test:
	go test -parallel 1 ./...
