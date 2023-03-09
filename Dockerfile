FROM golang:1.20

ENV MESSAGE_SERVICE_ENVIRONMENT=DEV
ENV MESSAGE_SERVICE_REPO_TYPE=POSTGRESQL
ENV MESSAGE_SERVICE_TIMEOUT=60
ENV MESSAGE_SERVICE_PORT=3333
ENV MESSAGE_SERVICE_PG_URL=postgres://postgres:postgres@yapyapyap-db:5432/postgres?sslmode=disable
ENV MESSAGE_SERVICE_TOKEN_SECRET=SECRET!

WORKDIR /go/src/github.com/stone1549/yapyapyap/message/
COPY . .

RUN apt-get update
RUN apt-get --assume-yes install libgeos-dev
RUN go mod tidy

CMD ["go", "run", "main.go"]

