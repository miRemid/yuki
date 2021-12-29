<h1 align="center">Welcome to Yuki 👋</h1>

[![build docker image](https://github.com/miRemid/yuki/actions/workflows/docker.yml/badge.svg)](https://github.com/miRemid/yuki/actions/workflows/docker.yml)

> A reverse proxy gateway for go-cqhttp

## 安装

### Release
在`Release`页面下载所需要的版本压缩包

### 自行编译

所需环境：
- Nodejs, yarn
- Go1.16
- make

```sh
git clone https://github.com/miRemid/yuki.git yuki
cd yuki
make web build-linux
# make web build-windows
cd release
```

## 使用方式

```sh
./yuki_linux_amd64
```
默认端口8080

```sh
❯ ./yuki_linux_amd64 -h
Usage of ./yuki_linux_amd64:
  -d    debug mode
  -p int
        server port (default 8080)
```

启动完成后，打开浏览器输入`http://127.0.0.1:8080`进入web管理界面

## API列表
Yuki内置了`swagger`文档，打开`http://127.0.0.1:8080/swagger/index.html即可

## Docker
现已支持Docker部署
```shell
docker pull kamir3mid/yuki:latest
docker run --name yuki -p 8080:8080 -v ${PWD}/data:/data -d yuki
```

## 

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_
