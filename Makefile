BIN := recanatid
ENTRY := cmd/recanatid/main.go
COVERFILE := cover.out

run:
	@go run $(ENTRY)

test:
	@go test -v ./...

build: test
	@go build -v -o $(BIN) $(ENTRY) \
	&& strip $(BIN)

coverage: # See golang.org/x/tools/cmd/cover
	@go test -coverprofile $(COVERFILE) ./... && \
	go tool cover -html=$(COVERFILE)