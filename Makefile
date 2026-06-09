# Build frontend and backend for release

.PHONY: all frontend backend clean

all: frontend backend

frontend:
	@echo "==> Building frontend in ./web/ui"
	cd ./web/ui && yarn build

backend:
	@echo "==> Building backend statically"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./vmcontroller .

clean:
	@echo "==> Cleaning build artifacts"
	rm -f ./vmcontroller
