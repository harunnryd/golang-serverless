build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/jukir_alpr_car jukir_alpr_car/main.go