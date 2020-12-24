build: test
	go build -o ./bin/mbserver mathbattle/cmd/mbserver
	go build -o ./bin/mbbot mathbattle/cmd/mbbot
	go build -o ./bin/adminclient mathbattle/cmd/adminclient
test:
	go test ./...