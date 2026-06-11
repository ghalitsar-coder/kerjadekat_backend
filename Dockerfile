# ==========================================
# Tahap 1: Builder
# ==========================================
FROM golang:alpine3.24 AS builder

ARG GOPROXY=https://goproxy.io,direct
ENV GOPROXY=$GOPROXY

WORKDIR /app

# Copy file dependency dan download modul
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi dari cmd/api/main.go
# CGO_ENABLED=0 memastikan binary yang dihasilkan bisa jalan mandiri di Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/kerjadekat-api ./cmd/api

# ==========================================
# Tahap 2: Final Image
# ==========================================
FROM alpine:latest

WORKDIR /app

# Hanya copy file binary hasil build dari tahap 1
COPY --from=builder /app/kerjadekat-api .

# (Opsional) Copy folder config jika aplikasi Anda butuh baca file config secara langsung
COPY --from=builder /app/config /app/config 

# Sesuaikan dengan port yang digunakan di main.go Anda (misal 8080)
EXPOSE 8085

# Jalankan aplikasi
CMD ["./kerjadekat-api"]