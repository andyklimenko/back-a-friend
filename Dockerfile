FROM golang:1.8.3

COPY . .

RUN go-wrapper download github.com/mattn/go-sqlite3
RUN go-wrapper install github.com/mattn/go-sqlite3

CMD ["go", "run", "src/back-a-friend.go"]
