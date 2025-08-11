# AI Agent Instructions

## Project Overview
OSC-to-MIDI bridge written in Go with Docker-based development. Provides bidirectional, ultra-low latency (1-2ms) conversion between OSC messages and JACK MIDI events.

## Essential Commands
```bash
make build           # Build the binary
make test            # Run unit tests (required before completion)
make fmt             # Format Go code (required before completion)
make vet             # Run Go vet linter (required before completion)
make integration-test # Run end-to-end validation
make dev             # Start development container
make run             # Run application in Docker
make clean           # Clean up containers
```

## Development Guidelines

### Prerequisites
- Read README.md first for user-facing overview
- All development happens inside Docker containers
- Docker required on macOS (JACK dependency)
- Use `gh` commands instead of fetch/curl for GitHub repos

### Code Quality Requirements
Before marking any task complete:
1. All unit tests must pass (`make test`)
2. Code must be formatted (`make fmt`) 
3. Linter must pass (`make vet`)
4. Build must succeed (`make build`)
5. Write minimal, simple, focused code

### Key Technical Details
- **Performance**: Ultra-low latency (1.33ms @ 48kHz), lock-free design, batch processing
- **Architecture**: Bidirectional OSC↔MIDI with separate queues (1024 MIDI events, 16 OSC messages)
- **Debugging**: Use `DEBUG=*` or specific modules (`DEBUG=bridge,main,handlers`)
- **Graceful shutdown**: SIGTERM/SIGINT handling implemented

## Development Workflows

### Docker Development
All commands run in containers with JACK dummy driver:
```bash
make dev                    # Interactive development
make run ARGS="--osc-port 8000"  # Run with custom args
```

### Testing
- **Unit tests**: `make test` (must pass for all changes)
- **Integration tests**: `make integration-test` (verifies OSC handling, JACK integration, graceful shutdown)

### Native Linux Build
```bash
sudo apt-get install jackd2 libjack-jackd2-dev
go mod download
go build -o osc-midi-bridge
jackd -d alsa -r 48000 -p 64    # or dummy driver for testing
DEBUG=* ./osc-midi-bridge       # with debug output
```

## Architecture Reference

### Event Flow
```
OSC→MIDI: UDP:9000 → OSC Server → Lock-free Queue → JACK Callback → MIDI Out
MIDI→OSC: MIDI In → JACK Callback → Buffered Channel → OSC Sender → UDP:8000
```

### Real-time Considerations
- JACK callback runs in real-time context
- Lock-free queue (OSC→MIDI) prevents priority inversion  
- Buffered channel (MIDI→OSC) avoids blocking RT thread
- Pre-allocated buffers minimize allocations

### Debug Namespaces
- `main` - Program flow
- `bridge` - JACK/MIDI operations  
- `handlers` - OSC message handling