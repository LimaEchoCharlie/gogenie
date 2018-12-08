# Gogenie 

Gogenie provides a simple GO interface to control an Energenie ENER002-2PI remote control socket via a Raspberry Pi 
controller board. Gogenie uses [periph.io](https://periph.io/) for GPIO control and has been tested on a Raspberry Pi 3. 

## Example
~~~go
import (
	"github.com/limaechocharlie/gogenie"
	"time"
)

// choose which plug to control
p := gogenie.PlugOne

// switch the plug on
p.On()

// the controller believes that the plug is now on
fmt.Printf(“The plug is on? %s\n”, p.IsOn())

// wait for a bit
time.Sleep(4 * time.Second)

// switch the plug off
p.Off()
~~~
## Development Build
By default, gogenie builds for a raspberry pi but a *dev* tag is available for development and testing away from the target system.
A *dev* build mocks the GPIO pins and logs any changes in pin state. 
To build the *dev* version, use `go build -tags 'dev’`.
