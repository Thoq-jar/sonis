FROM golang

WORKDIR /app

COPY . .

RUN go build -v .

EXPOSE 8080

CMD ["./sonis"]
