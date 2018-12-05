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
	PlugAll PlugID = iota
	PlugOne
	PlugTwo
)

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

// newPlug creates a new variable to control the plug with the supplied id
func NewPlug(id PlugID) *plug {
	p := &plug{
		id: id,
	}

	// initialise the plugs
	if err := initPlugs(); err != nil {
		log.Fatal(err)
	}

	// start with plug off
	p.Set(false)

	return p
}

// set sets the plug
func (p *plug) Set(on bool) error {
	// lock pins
	mutex.Lock()
	defer mutex.Unlock()

	// clear error
	clearPinError()

	// set d2-d1-d0 depending on which plug
	switch p.id {
	case PlugAll:
		// 011
		d2.off()
		d1.on()
		d0.on()
	case PlugOne:
		// 111
		d2.on()
		d1.on()
		d0.on()
	case PlugTwo:
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

// state returns the current status of the plug
func (p *plug) State() bool {
	return p.state
}
