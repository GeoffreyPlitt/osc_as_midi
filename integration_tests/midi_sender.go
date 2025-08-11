package main
import("time";"github.com/xthexder/go-jack")
type MidiEvent struct{status,note,velocity,delay int}
func main(){
	client,_:=jack.ClientOpen("midi_sender",jack.NoStartServer);defer client.Close()
	outPort:=client.PortRegister("out",jack.DEFAULT_MIDI_TYPE,jack.PortIsOutput,0)
	eventQueue:=make(chan[]byte,4)
	client.SetProcessCallback(func(nframes uint32)int{
		buffer:=outPort.MidiClearBuffer(nframes)
		for{select{case data:=<-eventQueue:outPort.MidiEventWrite(&jack.MidiData{Buffer:data,Time:0},buffer);default:return 0}}
	})
	client.Activate();time.Sleep(500*time.Millisecond)
	for _,event:=range[]MidiEvent{{0x90,67,64,250},{0x80,67,64,50},{0x90,79,64,250},{0x80,79,64,100}}{
		eventQueue<-[]byte{byte(event.status),byte(event.note),byte(event.velocity)}
		time.Sleep(time.Duration(event.delay)*time.Millisecond)
	}
}
