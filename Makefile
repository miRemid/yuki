.PHONY: build

RELEASE_DIR="release"
LINUX_AMD64_NAME="yuki_linux_amd64"
WINDOWS_AMD64_NAME="yuki_windows_x86-64.exe"


run:
	@go run *.go -d

build-linux:
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -a -ldflags \
	' -extldflags "-static"' \
	-o ${RELEASE_DIR}/${LINUX_AMD64_NAME}

build-windows:
	CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -a -ldflags \
	' -extldflags "-static"' \
	-o ${RELEASE_DIR}/${WINDOWS_AMD64_NAME}
