# df-tpstracker

A dragonfly library to track the TPS (ticks per second) of a world.

## Install

```sh
go get github.com/akmalfairuz/df-tpstracker
```

## Usage

```go
// import the library
import (
	...
	"github.com/akmalfairuz/df-tpstracker"
	...
)

var w *world.World

tracker := tpstracker.New(w)
go tracker.StartTracking()

// tracker.TPS(n) where n is the number of sample ticks to calculate the TPS.
// Max n is 600 (30 seconds).

// Get the current TPS
tps := tracker.TPS(1)

// Get the average TPS in the last 1 second (20 ticks)
tps := tracker.TPS(20)

// Don't forget to close the tracker when the world is closed.
func (wh worldHandler) HandleClose(tx *world.Tx) {
    tracker.Close()
}
```
