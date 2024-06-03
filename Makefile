BINARY_NAME=calendarbot
BINARY_DIR=bin

build:
	go build -o ${BINARY_DIR}/${BINARY_NAME} ./...
	GOARCH=arm64 GOOS=linux go build -o ${BINARY_DIR}/${BINARY_NAME}-linux ./...

test:
	go test ./...

test-integration:
	go test ./... -tags=integration

deploy: build
	scp -i ~/.ssh/personal_rsa bin/${BINARY_NAME}-linux idoberko2@raspberrypi.local:~/calendarbot/

run: build
	./${BINARY_DIR}/${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_DIR}/${BINARY_NAME}
	rm ${BINARY_DIR}/${BINARY_NAME}-linux
