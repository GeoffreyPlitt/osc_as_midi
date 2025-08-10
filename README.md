# Custom Go OSC-to-Virtual-MIDI Bridge

## Architecture Overview

A Go program that:
1. **Creates virtual JACK MIDI ports** using go-jack
2. **Listens for OSC messages** on configurable UDP port
3. **Translates OSC to MIDI** with flexible path patterns
4. **Emits MIDI events** through the virtual port with ultra-low latency (1.33ms)
5. **Logs debug info** using debuggo
6. **Gracefully cleans up** on shutdown

## Technical Stack

**Libraries Used:**
- `github.com/hypebeast/go-osc/osc` - Pure Go OSC implementation
- `github.com/xthexder/go-jack` - JACK Audio Connection Kit bindings
- `github.com/GeoffreyPlitt/debuggo` - Debug logging (controlled via DEBUG env var)
- System dependency: `jackd2` and `libjack-jackd2-dev` (for JACK)

## Performance

- **Ultra-low latency**: 1.33ms with 64 frame buffer @ 48kHz
- **Lock-free design**: Single producer/consumer event queue
- **Batch processing**: Up to 32 MIDI events per JACK cycle
- **No scheduling overhead**: Direct dispatch in real-time callback

## CLI Interface

```bash
# Default usage
./osc-midi-bridge

# With options
./osc-midi-bridge \
  --osc-port 8000 \
  --client-name "MyController" \
  --port-name "midi_out"

# List available MIDI ports
./osc-midi-bridge --list-ports
```

**CLI Flags:**
```
--osc-port      UDP port for OSC messages (default: 9000)
--client-name   JACK client name (default: "osc-midi-bridge")
--port-name     JACK MIDI output port name (default: "midi_out")
--list-ports    List available MIDI ports and exit
```

**Note:** JACK buffer size is controlled externally via the `jackd` command (e.g., `jackd -p 64`).

## OSC Message Format

**Supported paths:**
- `/midi/{channel}/note_on` - args: [note(int), velocity(int)]
- `/midi/{channel}/note_off` - args: [note(int), velocity(int)]

Where `{channel}` is 0-15 for MIDI channels 1-16.

**Example messages:**
- `/midi/0/note_on 60 127` - Middle C, channel 1, full velocity
- `/midi/9/note_on 36 100` - Kick drum, channel 10
- `/midi/0/note_off 60 0` - Middle C off, channel 1

## Development with Docker

The project uses Docker for cross-platform development, especially useful on macOS where JACK requires special setup. The Docker container runs JACK with a dummy driver for development.

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
DEBUG=* ./osc-midi-bridge
```

## Docker Usage

The Docker setup automatically handles JACK configuration:

```bash
# Build and run with Docker
make run

# This starts JACK with dummy driver and runs the bridge
# Equivalent to:
# docker run ... bash -c "jackd -d dummy -r 48000 -p 64 & sleep 1 && ./osc-midi-bridge"
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

# Custom port and names
./osc-midi-bridge --osc-port 8000 --client-name "OSC Controller" --port-name "osc_out"

# Ultra-low latency (0.67ms @ 48kHz)
jackd -d dummy -r 48000 -p 32 &
./osc-midi-bridge

# List available MIDI ports
./osc-midi-bridge --list-ports
```

## Connecting to Other Applications

Once running, the bridge creates a JACK MIDI output port that can be connected to any JACK-aware application:

```bash
# List all JACK connections
jack_lsp -c

# Connect to a software synthesizer
jack_connect midi_out yoshimi:midi_in

# Monitor MIDI output
jack_midi_dump midi_out
```

## Sending OSC Messages

You can send OSC messages using any OSC client. For testing, use `oscsend`:

```bash
# Send middle C note on
oscsend localhost 9000 /midi/0/note_on ii 60 127

# Send middle C note off
oscsend localhost 9000 /midi/0/note_off ii 60 0

# Send a chord
oscsend localhost 9000 /midi/0/note_on ii 60 100 &
oscsend localhost 9000 /midi/0/note_on ii 64 100 &
oscsend localhost 9000 /midi/0/note_on ii 67 100 &
```

## Architecture Details

### Event Flow
```
OSC Client → UDP:9000 → OSC Server → Event Queue → JACK Process Callback → JACK MIDI Port
                                          ↓
                                    (Lock-free FIFO)
                                    1024 event buffer
```

### Real-time Considerations
- JACK process callback runs in real-time context
- Lock-free event queue prevents priority inversion
- Batch processing limits per-cycle CPU usage
- Pre-allocated buffers minimize allocations

### Debug Logging
The bridge uses hierarchical debug namespaces:
- `main` - Main program flow
- `bridge` - JACK client and MIDI operations
- `handlers` - OSC message handling

Enable specific namespaces with the DEBUG environment variable.