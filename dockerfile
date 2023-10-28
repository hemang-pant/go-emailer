FROM golang:1.21.3

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY template.html .

RUN go get
RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]