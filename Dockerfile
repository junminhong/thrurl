FROM golang
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
ADD . /app
RUN cd /app/cmd && go build -o /app/docker-thrurl
EXPOSE 9020
CMD [ "/app/docker-thrurl" ]