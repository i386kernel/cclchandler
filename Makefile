
build:
	GOOS=linux GOARCH=amd64 go build -o cclcmgr-v1.0.0-beta-linux-amd64 .
	chmod +x cclcmgr-v1.0.0-beta-linux-amd64
