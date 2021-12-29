# build react
FROM node:14.17.5-alpine as node

ADD . /yuki
WORKDIR /yuki/web

RUN yarn && yarn build

# build golang
FROM tetafro/golang-gcc:1.16-alpine as golang

COPY --from=node /yuki /yuki
WORKDIR /yuki

RUN CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -ldflags \
	' -extldflags "-static"' \
	-o yuki_linux_amd64

# finish build
FROM alpine:3.14

COPY --from=golang /yuki/yuki_linux_amd64 /yuki_linux_amd64
COPY --from=golang /yuki/web/dist /web/dist
COPY --from=golang /yuki/docs /docs

WORKDIR /

CMD [ "sh", "-c", "/yuki_linux_amd64" ]