# Stage 1: Build the Go binary
FROM golang:1.24@sha256:b6652731d1d5622f85509f72942ce2a344e3bf6dd6793b2e462cc5fb3126b566 AS builder

WORKDIR /app

COPY app/go.mod app/go.sum ./
RUN go mod download

COPY app/ ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o filestore-backup main.go

# Stage 2: Create a lightweight image
FROM alpine:3.21.3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

# Create minimal, non-root user
RUN addgroup -S filestoregroup && adduser -S filestoreuser -G filestoregroup

WORKDIR /app

COPY --from=builder /app/filestore-backup .

RUN chmod +x /app/filestore-backup
RUN chown -R filestoreuser:filestoregroup /app

USER filestoreuser

ENTRYPOINT ["./filestore-backup"]
