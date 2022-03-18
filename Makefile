build:
		GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o thomas

release: build
		upx -9 thomas