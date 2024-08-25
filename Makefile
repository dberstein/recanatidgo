BIN := recanatid
ENTRY := cmd/recanatid/main.go

run:
	@go run $(ENTRY)

build:
	@go build -v -o $(BIN) $(ENTRY) \
	&& strip $(BIN)