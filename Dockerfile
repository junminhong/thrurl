# 用alpine版本創建出來的docker容量比較小，不會動不動就破GB起跳
FROM golang:1.17.6-alpine
# 因為golang下載套件需要使用到git
RUN apk add --no-cache git
# RUN go get github.com/lib/pq
# 設置工作目錄
WORKDIR /app
# 複製go.mod and go.sum
COPY go.mod go.sum ./
# 將所有檔案放進去app
ADD . /app
# build exec 下載go mod package
RUN go mod download
# 定義暴露端口
EXPOSE 9220
