
build:
	GOOS=linux GOARCH=amd64 go build -o customcerthandler .
	chmod +x customcerthandler
