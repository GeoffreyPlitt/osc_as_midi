# Instructions for AI Agents

## Getting Started
1. **Read the README.md first** - It contains the project overview, architecture, and development workflow
2. **Understand the project structure** - This is an OSC-to-MIDI bridge written in Go with Docker-based development


## Remember
- You cannot fetch/curl to read github repos, use the `gh` command instead. Instead of fetching or cloning the repo, use `gh api repos/.../.../contents` to view contents, etc.
- Write the minimum code necessary to achieve the goal
- Keep implementations simple and focused
- Before marking task as complete:
- Run unit tests: `make test`
- Format code: `make fmt`
- Run linter: `make vet`
- Build the binary: `make build`
- All tests must pass before considering the work done
- Testing:
  - Unit tests: `make test` - Must pass for all code changes
  - Integration tests: `make integration-test` - For end-to-end validation

## Key Commands
```bash
make build           # Build the binary
make test            # Run unit tests
make fmt             # Format Go code
make vet             # Run Go vet linter
make integration-test # Run integration tests
make dev             # Start development container
```

### Important Notes
- All development happens inside Docker containers (see Makefile)
- The project uses JACK for MIDI on Linux, hence Docker is required on macOS
- Signal handling (SIGTERM/SIGINT) is implemented for graceful shutdown

# Performance

- **Ultra-low latency**: 1.33ms with 64 frame buffer @ 48kHz
- **Lock-free design**: Single producer/consumer event queue
- **Batch processing**: Up to 32 MIDI events per JACK cycle
- **No scheduling overhead**: Direct dispatch in real-time callback

# Development Workflow

## Docker Development

The Docker container runs JACK with a dummy driver for development:

```bash
# Build the binary
make build

# Format code
make fmt

# Run go vet
make vet

# Run unit tests
make test

# Start interactive development container
make dev

# Run the application
make run
make run ARGS="--osc-port 8000"

# Clean up
make clean
```

## Testing

### Unit Tests

Run the Go unit tests with coverage:

```bash
make test
```

### Integration Tests

Integration tests verify the OSC message handling and signal processing:

```bash
make integration-test
```

The integration test script:
- Verifies the bridge starts correctly
- Sends test OSC messages using `oscsend`
- Tests graceful shutdown with SIGTERM
- Validates the test tooling is properly configured

## Build Instructions (Linux)

```bash
# Install JACK development headers
sudo apt-get install jackd2 libjack-jackd2-dev

# Get dependencies
go mod download

# Build
go build -o osc-midi-bridge

# Start JACK (if not already running)
jackd -d alsa -r 48000 -p 64  # For real hardware
# OR
jackd -d dummy -r 48000 -p 64  # For testing without hardware

# Run
./osc-midi-bridge

# Run with debug output
DEBUG=* ./osc-midi-bridge
```

## Example Usage Scenarios

```bash
# Basic usage with defaults
./osc-midi-bridge

# With debug output enabled
DEBUG=* ./osc-midi-bridge

# Enable debug for specific modules
DEBUG=handlers ./osc-midi-bridge
DEBUG=bridge,main ./osc-midi-bridge

# Custom ports and names
./osc-midi-bridge --osc-port 8000 --osc-target-host 192.168.1.100 --osc-target-port 9000 --client-name "OSC Controller" --port-name "osc_out"

# Ultra-low latency (0.67ms @ 48kHz)
jackd -d dummy -r 48000 -p 32 &
./osc-midi-bridge

# List available MIDI ports
./osc-midi-bridge --list-ports
```

# Architecture Details

## Event Flow (Bidirectional)
```
OSC → MIDI Direction:
OSC Client → UDP:9000 → OSC Server → Event Queue → JACK Process Callback → JACK MIDI Out Port
                                          ↓
                                    (Lock-free FIFO)
                                    1024 event buffer

MIDI → OSC Direction:  
JACK MIDI In Port → JACK Process Callback → OSC Queue → OSC Sender → UDP:8000 → OSC Client
                                                ↓
                                        (Buffered channel)
                                         16 message buffer
```

## Real-time Considerations
- JACK process callback runs in real-time context
- Lock-free event queue (OSC→MIDI) prevents priority inversion
- Buffered channel (MIDI→OSC) prevents blocking RT thread
- Batch processing limits per-cycle CPU usage
- Pre-allocated buffers minimize allocations
- OSC sending happens in separate goroutine to maintain real-time safety

## Debug Logging
The bridge uses hierarchical debug namespaces:
- `main` - Main program flow
- `bridge` - JACK client and MIDI operations
- `handlers` - OSC message handling

Enable specific namespaces with the DEBUG environment variable.