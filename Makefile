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

cover: # See golang.org/x/tools/cmd/cover
	@go test -coverprofile $(COVERFILE) ./... && \
	go tool cover -html=$(COVERFILE)

cover/func: test
	@go tool cover -func=$(COVERFILE)

cover/pct: test
	@$(MAKE) cover/func | grep total: | awk '{print $3}'