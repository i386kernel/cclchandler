
build:
	GOOS=linux GOARCH=amd64 go build -o cclchandler-v1.0.0-beta-linux-amd64 .
	chmod +x cclchandler-v1.0.0-beta-linux-amd64
