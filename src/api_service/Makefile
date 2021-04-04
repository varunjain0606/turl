PROGNAME = tiny_url_api
MAIN_GO = server/main/main.go
SRC = $(MAIN_GO)

all: build

build: $(MAIN_GO) $(SRC)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(PROGNAME) $(SRC)

clean:
	go clean ./...
	rm -f $(PROGNAME) server/restapi/doc.go server/restapi/embedded_spec.go server/restapi/server.go
	rm -rf server/cmd server/models server/restapi/operations

$(MAIN_GO): swagger.yaml
	swagger generate server -t server -f swagger.yaml -A tinyurlapi

linux: $(SRC)
	CGO_ENABLED=0 GOOS=linux go build -o $(PROGNAME) $(SRC)

image: linux
	docker build -t $(PROGNAME) -f Dockerfile .

publish: image publish-no-deps

publish-no-deps:
	$(call docker_push_images,$(PROGNAME))

.PHONY: all build clean linux image publish publish-no-deps