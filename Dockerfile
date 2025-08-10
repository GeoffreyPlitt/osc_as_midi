FROM golang:1.22-bookworm

# Install ALSA development libraries and tools
RUN apt-get update && apt-get install -y \
    libasound2-dev \
    alsa-utils \
    liblo-tools \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod files (when they exist)
COPY go.mod* go.sum* ./

# Download dependencies (if go.mod exists)
RUN if [ -f go.mod ]; then go mod download; fi

# Copy source code
COPY . .

# Build the application
RUN if [ -f go.mod ]; then \
        go mod tidy && \
        go build -o osc-midi-bridge; \
    fi

# Default command (use DEBUG env var for debug output)
CMD ["./osc-midi-bridge"]