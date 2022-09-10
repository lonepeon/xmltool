BIN_PATH := target/xmltool

.PHONY: build
build: $(BIN_PATH)-linux-amd64 $(BIN_PATH)-windows-amd64

$(BIN_PATH):
	go build -o $@

$(BIN_PATH)-linux-amd64: export GOOS=linux
$(BIN_PATH)-linux-amd64: export GOARCH=amd64
$(BIN_PATH)-linux-amd64:
	go build -o $@

$(BIN_PATH)-windows-amd64: export GOOS=windows
$(BIN_PATH)-windows-amd64: export GOARCH=amd64
$(BIN_PATH)-windows-amd64:
	go build -o $@

.PHONY: clean
clean:
	rm -rf $(dir $(BIN_PATH)) testfiles
