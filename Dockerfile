# syntax=docker/dockerfile:1

# ===================================================================================
# === Stage 1:Builder container =====================================================
# ===================================================================================
FROM golang:1.18-alpine AS build

WORKDIR /build

# Fetch modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy Go source files
COPY config/ ./config/
COPY cmd/ ./cmd/
COPY pkg/ ./pkg/
COPY internal/ ./internal/

RUN GOOS=linux \
    go build \
    -o server github.com/kapitan123/telegrofler/cmd/api

# ===================================================================================
# === Stage 2: Create a lightweight container =======================================
# ===================================================================================
FROM alpine

WORKDIR /app

COPY --from=build /build/server . 

EXPOSE 9001

CMD [ "./server"]