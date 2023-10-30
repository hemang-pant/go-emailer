FROM golang:1.21.3

WORKDIR /app

ENV PORT=9000

COPY go.mod .
COPY go.sum .
COPY main.go .
COPY template.html .
COPY .env .

RUN go get
RUN go build -o bin .

EXPOSE 9000

ENTRYPOINT [ "/app/bin", "--port", "9000:9000"]