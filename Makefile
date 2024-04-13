run: build
	@./bin/goredis --listenAddr :5002
build:
	@go build -o bin/goredis .
watch: 
	@air
