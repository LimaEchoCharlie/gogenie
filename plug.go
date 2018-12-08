package gogenie

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// mutex to protect pin writes
var mutex = &sync.Mutex{}

// plug ids
//go:generate stringer -type=PlugID
type PlugID int

const (
	all PlugID = iota
	one
	two
)

var (
	PlugAll = &plug{id: all}
	PlugOne = &plug{id: one}
	PlugTwo = &plug{id: two}
)

func init() {
	// initialise the plugs
	if err := initPlugs(); err != nil {
		log.Fatal(err)
	}
}

// initPlugs initialises the pins used to communicate with the plugs
func initPlugs() (err error) {
	// lock mutex
	mutex.Lock()
	defer mutex.Unlock()

	// clear error
	clearPinError()

	// initialise periph
	if err := initHAL(); err != nil {
		return err
	}

	// set encoder to 0000
	d3.off()
	d2.off()
	d1.off()
	d0.off()

	// diable modulator
	enable.off()

	// set modulator to ASK
	mode.off()
	return lastPinError()
}

type plug struct {
	id    PlugID
	timer *time.Timer
	state bool
}

// stringer function for a Plug
func (p plug) String() string {
	return fmt.Sprintf("Plug %s", p.id)
}

// set sets the plug to the given value; true = on, false = off
func (p *plug) set(on bool) error {
	// lock pins
	mutex.Lock()
	defer mutex.Unlock()

	// clear error
	clearPinError()

	// set d2-d1-d0 depending on which plug
	switch p.id {
	case all:
		// 011
		d2.off()
		d1.on()
		d0.on()
	case one:
		// 111
		d2.on()
		d1.on()
		d0.on()
	case two:
		// 110
		d2.on()
		d1.on()
		d0.off()
	default:
		// not recognised, return error
		return fmt.Errorf("%d is not a valid plug id", p.id)
	}

	// set d3 depending on on/off
	if on {
		d3.on()
	} else {
		d3.off()
	}

	// allow the encoder to settle
	time.Sleep(100 * time.Millisecond)

	// enable the modulator
	enable.on()
	// pause
	time.Sleep(250 * time.Millisecond)
	// disable the modulator
	enable.off()

	err := lastPinError()
	if err == nil {
		p.state = on
	}
	return err
}

// On sends a switch on message to the plug
func (p *plug) On() error {
	return p.set(true)
}

// Off sends a switch off message to the plug
func (p *plug) Off() error {
	return p.set(false)
}

// IsOn returns true if the current state indicates that the plug is on
// Note: this is the state that the controller believes the plug to be in since the communication is only one-way
func (p *plug) IsOn() bool {
	return p.state
}
