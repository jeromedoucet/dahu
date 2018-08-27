APP_NAME = dahu-server
DB_NAME = dahu
VERSION = $(shell git rev-parse --short HEAD)
LD_FLAGS = -ldflags "-X main.version=$(VERSION)"

all: $(APP_NAME)

$(APP_NAME): dep build

dep:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	dep ensure

build:
	@echo "Building application"
	go build $(LD_FLAGS) -o $(APP_NAME)

test: build
	@echo "Running tests..."
	go test -p 1 -coverpkg=./... -coverprofile=coverage.out ./...

showCov: test
	go tool cover -html=coverage.out

clean:
	rm -f $(APP_NAME)
	rm -f $(DB_NAME)

re: clean all
