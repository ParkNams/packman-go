FROM golang

WORKDIR /app

COPY . /app

RUN go mod download

CMD [ "go","run","main.go" ]