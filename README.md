# Custom Go OSC-to-Virtual-MIDI Bridge

[![Go Reference](https://pkg.go.dev/badge/github.com/GeoffreyPlitt/osc_as_midi.svg)](https://pkg.go.dev/github.com/GeoffreyPlitt/osc_as_midi)
[![Go Report Card](https://goreportcard.com/badge/github.com/GeoffreyPlitt/osc_as_midi)](https://goreportcard.com/report/github.com/GeoffreyPlitt/osc_as_midi)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![CI](https://github.com/GeoffreyPlitt/osc_as_midi/workflows/CI/badge.svg)](https://github.com/GeoffreyPlitt/osc_as_midi/actions)

## Architecture Overview

A **bidirectional** Go program that:
1. **EMULATES a VIRTUAL Jack MIDI device** (input and output) using go-jack
2. **OSC → MIDI**: Listens for OSC messages on configurable UDP port, and emits them as Note events from the virtual device.
3. **MIDI → OSC**: Reads incoming MIDI messages sent to the virtual MIDI device, and emits them as OSC messages.
5. **Emits MIDI/OSC events** Low latency (1-2ms)
6. **Logs debug info** using debuggo
7. **Gracefully cleans up** on shutdown

## Technical Stack

**Libraries Used:**
- `github.com/hypebeast/go-osc/osc` - Pure Go OSC implementation
- `github.com/xthexder/go-jack` - JACK Audio Connection Kit bindings
- `github.com/GeoffreyPlitt/debuggo` - Debug logging (controlled via DEBUG env var)
- System dependency: `jackd2` and `libjack-jackd2-dev` (for JACK)

## CLI Interface

```bash
./osc-midi-bridge \
  --osc-port 8000 \
  --osc-target-host 192.168.1.100 \
  --osc-target-port 9000 \
  --client-name "OSC Controller" \
  --port-name "osc_out"
```

**CLI Flags:**
```
--osc-port         UDP port for incoming OSC messages (default: 9000)
--osc-target-host  Target host for outgoing OSC messages (default: "localhost")
--osc-target-port  Target port for outgoing OSC messages (default: 8000)
--client-name      JACK client name (default: "osc-midi-bridge")
--port-name        JACK MIDI output port name (default: "midi_out")
--list-ports       List available MIDI ports and exit
```

**Note:** JACK buffer size is controlled externally via the `jackd` command (e.g., `jackd -p 64`).

## Message Format

**OSC Paths (Bidirectional):**
- `/midi/{channel}/note_on` - args: [note(int), velocity(int)]
- `/midi/{channel}/note_off` - args: [note(int), velocity(int)]

Where `{channel}` is 0-15 for MIDI channels 1-16.

**Example OSC Messages:**
- `/midi/0/note_on 60 127` - Middle C, channel 1, full velocity
- `/midi/9/note_on 36 100` - Kick drum, channel 10
- `/midi/0/note_off 60 0` - Middle C off, channel 1

**Bidirectional Flow:**
- **Incoming OSC** → **Outgoing MIDI**: Messages received on `--osc-port` (default 9000) are converted to MIDI and sent via JACK `midi_out` port
- **Incoming MIDI** → **Outgoing OSC**: MIDI events received via JACK `midi_in` port are converted to OSC and sent to `--osc-target-host:--osc-target-port` (default localhost:8000)

## Development

The project uses Docker for cross-platform development, especially useful on macOS where JACK requires special setup. Run `make test` for unit tests and `make integration-test` for end-to-end validation. See CLAUDE.md for detailed development workflow.