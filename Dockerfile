FROM golang:1.20

WORKDIR /app

COPY . .

RUN go get -d -v ./...
RUN go build -o bin/movie_ticket_back_api.exe ./cmd
RUN [ "chmod", "+x", "/app/bin/movie_ticket_back_api.exe"]

ENTRYPOINT [ "./bin/movie_ticket_back_api.exe" ]

EXPOSE 3100
