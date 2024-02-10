# Executable name
EXC_NAME=snow

# Directory location
CMD_PATH=./cmd/


build:
	go build -ldflags="-linkmode external -extldflags 'app.res -static'" -o ${EXC_NAME}.exe ${CMD_PATH}

build-linux:
	GOOS=linux GOARCH=amd64 go build -o snow ./cmd/

test:
	go test -v ./...


run:
	go run ./cmd

clean:
	go clean
	rm .\${EXC_NAME}.exe

uac-emb:
	rsrc -manifest app.manifest -o rsrc.syso

# -s -w strips debugging info disables DWARF symbol table generation (reduced binary size)
strip-build:
	go build -ldflags="-s -w -linkmode external -extldflags 'app.res -static'" -o ${EXC_NAME}.exe ${CMD_PATH}