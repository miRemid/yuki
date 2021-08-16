.PHONY: build

LINUX_AMD64_NAME="yuki_linux_amd64"
WINDOWS_AMD64_NAME="yuki_windows_x86-64.exe"

WORKDIR=$(shell pwd)
RELEASE_DIR=${WORKDIR}/release
DATA_DIR=${WORKDIR}/data

clean:
	rm -rf ${DATA_DIR} ${RELEASE_DIR}

run:
	@go run *.go -d

pre:
	mkdir -p ${RELEASE_DIR}/web

web: pre
	cd web && yarn && yarn build && cp -r dist ${RELEASE_DIR}/web

build-linux:
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -a -ldflags \
	' -extldflags "-static"' \
	-o ${RELEASE_DIR}/${LINUX_AMD64_NAME}

build-windows:
	CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -a -ldflags \
	' -extldflags "-static"' \
	-o ${RELEASE_DIR}/${WINDOWS_AMD64_NAME}
