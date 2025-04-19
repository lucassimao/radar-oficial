FROM golang:1.24.2

# Install dependencies
RUN apt-get update && apt-get install -y \
    wget curl gnupg ca-certificates \
    git build-essential gcc \
    libglib2.0-0 libnss3 libgconf-2-4 libfontconfig1 libxss1 libasound2 libxtst6 \
    libatk-bridge2.0-0 libgtk-3-0 \
    software-properties-common

# Install Chrome (headless-compatible)
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
RUN apt install -y ./google-chrome-stable_current_amd64.deb

# Set working directory
WORKDIR /app

# Copy only go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies first (cached unless go.mod/sum changes)
RUN go mod download

# Copy code
COPY . .

# Build worker binary
RUN go build -o bin/worker ./cmd/worker

# Entrypoint
CMD ["./bin/worker"]
