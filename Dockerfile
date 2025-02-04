FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh
#RUN go mod tidy

COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY migrations/ ./migrations/

ENV GOOSE_DRIVER=postgres
ENV GOOSE_DBSTRING=postgres://postgres:somepostgres@db:5432/postgres
ENV GOOSE_MIGRATION_DIR=./migrations

RUN goose up

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /test-iq ./cmd/main.go

EXPOSE 8080

# Run
CMD ["/test-iq"]