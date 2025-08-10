FROM golang:1.22-bookworm

# Install JACK development libraries and tools
RUN apt-get update && apt-get install -y \
    jackd2 \
    libjack-jackd2-dev \
    liblo-tools \
    jack-tools \
    bsdmainutils \
    python3 \
    python3-pip \
    && pip3 install --break-system-packages honcho \
    && rm -rf /var/lib/apt/lists/*

# Set JACK environment variable to prevent audio reservation issues
ENV JACK_NO_AUDIO_RESERVATION=1

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