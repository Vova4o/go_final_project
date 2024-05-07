FROM golang:1.22

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/todo-app ./cmd/app

EXPOSE 7540

CMD ["/app/todo-app"]