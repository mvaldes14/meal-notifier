FROM golang:1.25.1

WORKDIR /app

COPY . /app

RUN go build -o meal-notifier .

ENTRYPOINT ["./meal-notifier"]
