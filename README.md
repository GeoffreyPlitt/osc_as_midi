# Custom Go OSC-to-Virtual-MIDI Bridge

## Architecture Overview

A Go program that:
1. **Creates virtual ALSA MIDI ports** using gomidi/rtmididrv
2. **Listens for OSC messages** on configurable UDP port
3. **Translates OSC to MIDI** with flexible path patterns
4. **Emits MIDI events** through the virtual port
5. **Logs debug info** using debuggo
6. **Gracefully cleans up** on shutdown

## Technical Stack

**Libraries Used:**
- `github.com/hypebeast/go-osc/osc` - Pure Go OSC implementation
- `gitlab.com/gomidi/midi/v2` - MIDI library
- `gitlab.com/gomidi/midi/v2/drivers/rtmididrv` - Virtual port creation
- `github.com/GeoffreyPlitt/debuggo` - Debug logging
- System dependency: `libasound2-dev` (for ALSA)

## CLI Interface

```bash
# Default usage
./osc-midi-bridge

# With options
./osc-midi-bridge \
  --osc-port 8000 \
  --midi-name "MyController" \
  --osc-pattern "/control/{channel}/{type}" \
  --debug

# List available MIDI ports
./osc-midi-bridge --list-ports
```

**CLI Flags:**
```go
flag.IntVar(&oscPort, "osc-port", 9000, "UDP port to listen for OSC messages")
flag.StringVar(&midiName, "midi-name", "Virtual MIDI", "Name of virtual MIDI port")
flag.BoolVar(&listPorts, "list-ports", false, "List available MIDI ports and exit")
```

## OSC Message Format
- Hardcoded format: `/midi/{channel}/{type}` 
- We can expand to make it configurable in the future.

**Message Types:**
- `/midi/{channel}/note_on` - args: [note(int), velocity(int)]
- `/midi/{channel}/note_off` - args: [note(int), velocity(int)]

## Development with Docker

For cross-platform development (especially on macOS), use the provided Makefile:

```bash
# Build the binary
make build

# Run tests
make test

# Format code
make fmt

# Run go vet
make vet

# Start interactive development container
make dev

# Run the application
make run
make run ARGS="--debug --osc-port 8000"

# Clean up
make clean
```

## Build Instructions (Linux)

```bash
# Install ALSA dev headers
sudo apt-get install libasound2-dev

# Get dependencies
go mod download

# Build
go build -o osc-midi-bridge

# Run
./osc-midi-bridge --debug
```

## Example Usage Scenarios

```bash
# Basic usage with defaults
./osc-midi-bridge

# Custom port and name
./osc-midi-bridge --osc-port 8000 --midi-name "OSC Controller"

# List available MIDI ports
./osc-midi-bridge --list-ports
```
