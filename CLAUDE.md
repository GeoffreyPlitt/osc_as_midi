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
- The project uses ALSA for MIDI on Linux, hence Docker is required on macOS
- Signal handling (SIGTERM/SIGINT) is implemented for graceful shutdown
- Debug logging uses the debuggo library with DEBUG environment variable (not --debug flag)
  - Enable all debug output: `DEBUG=*`
  - Enable specific modules: `DEBUG=osc-midi-bridge:*`
- You cannot fetch/curl to read github repos, use the `gh` command instead