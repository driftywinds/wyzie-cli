APP     := wyzie-subs
OUT     := dist
FLAGS   := -ldflags="-s -w"

.PHONY: build clean run

build:
	mkdir -p $(OUT)
	GOOS=linux   GOARCH=amd64  go build $(FLAGS) -o $(OUT)/$(APP)-linux-amd64
	GOOS=linux   GOARCH=arm64  go build $(FLAGS) -o $(OUT)/$(APP)-linux-arm64
	GOOS=darwin  GOARCH=amd64  go build $(FLAGS) -o $(OUT)/$(APP)-macos-amd64
	GOOS=darwin  GOARCH=arm64  go build $(FLAGS) -o $(OUT)/$(APP)-macos-arm64
	GOOS=windows GOARCH=amd64  go build $(FLAGS) -o $(OUT)/$(APP)-windows-amd64.exe
	@echo ""
	@echo "✓ All binaries written to ./$(OUT)/"
	@ls -lh $(OUT)/

run:
	go run .

clean:
	rm -rf $(OUT)
