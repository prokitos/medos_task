APP?=medos.exe
CONTNAME?=medos-container
VERSION?=0.0.1



currentDepoly: deployDocker



deployBase: buildBase test fmt vet 
deployDocker: buildBase test fmt vet docker




buildBase:
	go build -o ${APP} cmd/main.go
test:
	go test ./...
fmt:
	go fmt ./...
vet:
	go vet ./...
docker:
	docker build -t ${CONTNAME}:${VERSION} .
	docker compose -f docker-compose.yml up