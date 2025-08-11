package main

import (
	"fmt"
	"os"
	"time"

	"github.com/xthexder/go-jack"
)

func main() {
	// Connect to JACK
	client, status := jack.ClientOpen("midi_sender", jack.NoStartServer)
	if status != 0 {
		fmt.Printf("Cannot connect to JACK: %s\n", jack.StrError(status))
		os.Exit(1)
	}
	defer client.Close()

	// Create MIDI output port
	outPort := client.PortRegister("out", jack.DEFAULT_MIDI_TYPE, jack.PortIsOutput, 0)
	if outPort == nil {
		fmt.Println("Failed to create MIDI output port")
		os.Exit(1)
	}

	// Event queue to send MIDI events
	eventQueue := make(chan []byte, 32)

	// Set up process callback
	if code := client.SetProcessCallback(func(nframes uint32) int {
		buffer := outPort.MidiClearBuffer(nframes)

		// Send all queued events
		for {
			select {
			case midiData := <-eventQueue:
				midiEvent := &jack.MidiData{
					Buffer: midiData,
					Time:   0, // Send immediately
				}
				if err := outPort.MidiEventWrite(midiEvent, buffer); err != 0 {
					fmt.Printf("Failed to write MIDI event: %v\n", err)
				}
			default:
				return 0 // No more events to process
			}
		}
	}); code != 0 {
		fmt.Printf("Failed to set process callback: %s\n", jack.StrError(code))
		os.Exit(1)
	}

	// Activate the client
	if code := client.Activate(); code != 0 {
		fmt.Printf("Failed to activate client: %s\n", jack.StrError(code))
		os.Exit(1)
	}

	fmt.Println("MIDI sender activated")

	// Wait longer to ensure connection is established
	time.Sleep(500 * time.Millisecond)

	// Send deterministic sequence:
	// Note 67 on -> wait 0.25s -> Note 67 off -> Note 79 on -> wait 0.25s -> Note 79 off

	fmt.Println("Sending Note 67 ON")
	// Note 67 ON
	eventQueue <- []byte{0x90, 67, 64} // Note on channel 0, note 67, velocity 64
	time.Sleep(250 * time.Millisecond)

	fmt.Println("Sending Note 67 OFF")
	// Note 67 OFF
	eventQueue <- []byte{0x80, 67, 64} // Note off channel 0, note 67
	time.Sleep(50 * time.Millisecond)

	fmt.Println("Sending Note 79 ON")
	// Note 79 ON
	eventQueue <- []byte{0x90, 79, 64} // Note on channel 0, note 79, velocity 64
	time.Sleep(250 * time.Millisecond)

	fmt.Println("Sending Note 79 OFF")
	// Note 79 OFF
	eventQueue <- []byte{0x80, 79, 64} // Note off channel 0, note 79
	time.Sleep(100 * time.Millisecond)

	fmt.Println("MIDI sequence complete")
}
