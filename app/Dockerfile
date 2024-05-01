FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod tidy

RUN go build cmd/gophermart/main.go

CMD ["go", "run", "main"]

#RUN go get github.com/githubnemo/CompileDaemon

#ENV PATH=$PATH:/go/bin

#ENTRYPOINT CompileDaemon --build="go build cmd/gophermart/main.go"