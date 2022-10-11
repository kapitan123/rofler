# syntax=docker/dockerfile:1

# ===================================================================================
# === Stage 1:Builder container =====================================================
# ===================================================================================
FROM golang:1.19-alpine AS builder

WORKDIR /build

# Fetch modules
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy Go source files
COPY config/ ./config/
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN GOOS=linux \
    go build \
    -o server github.com/kapitan123/telegrofler/cmd/api

# ===================================================================================
# === Stage 2: Create a lightweight container =======================================
# ===================================================================================

# AK TODO transfer to scratch
FROM alpine 

WORKDIR /app

COPY --from=builder /build/server . 

EXPOSE 9001

CMD [ "./server"]