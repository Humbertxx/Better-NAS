# STAGE 1: Build the application
FROM golang:1.22 AS builder
WORKDIR /app

# Copy dependency files and download
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code and compile a statically linked binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

# STAGE 2: Run the application
FROM scratch
WORKDIR /app

# Copy ONLY the binary from the builder stage
COPY --from=builder /app/myapp .

# Run the executable
ENTRYPOINT ["./myapp"]