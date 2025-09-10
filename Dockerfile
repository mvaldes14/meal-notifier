FROM golang:1.21.5-alpine

WORKDIR /app

COPY . /app

RUN go build -o meal-notifier .

ENTRYPOINT ["./meal-notifier"]
