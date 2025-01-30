package tpstracker

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/world"
	"sync"
	"time"
)

// TPSTracker tracks TPS with accurate averages.
type TPSTracker struct {
	w             *world.World
	tickDurations []float64
	closing       chan struct{}
	wg            sync.WaitGroup
	once          sync.Once
}

// New initializes a new TPS tracker.
func New(w *world.World) *TPSTracker {
	td := make([]float64, 0, 600)
	for i := 0; i < 600; i++ {
		td = append(td, 0.05)
	}
	return &TPSTracker{w: w, tickDurations: td}
}

// StartTracking runs the TPS monitoring loop.
func (t *TPSTracker) StartTracking() {
	tc := time.NewTicker(time.Second / 20) // 20 TPS
	t.wg.Add(1)
	defer tc.Stop()

	for {
		select {
		case <-tc.C:
			t.measureTickDuration()
		case <-t.closing:
			t.wg.Done()
			return
		}
	}
}

// measureTickDuration measures the duration of a tick and stores it in the tracker.
func (t *TPSTracker) measureTickDuration() {
	start := time.Now()

	<-t.w.Exec(func(tx *world.Tx) {})

	elapsed := time.Since(start).Seconds()

	if elapsed < 0.05 {
		elapsed = 0.05
	}

	t.tickDurations = append(t.tickDurations, elapsed)
	if len(t.tickDurations) > 600 {
		t.tickDurations = t.tickDurations[1:]
	}
}

// TPS returns the TPS for the last n samples.
func (t *TPSTracker) TPS(samples int) float64 {
	sum := 0.0
	for _, d := range t.tickDurations[len(t.tickDurations)-samples:] {
		sum += d
	}
	return float64(samples) / sum
}

// PrintTPS prints the TPS for the last 20, 100, 200, and 600 samples.
func (t *TPSTracker) PrintTPS() {
	fmt.Printf("TPS: [1s: %.2f] [5s: %.2f] [10s: %.2f] [30s: %.2f]\n", t.TPS(20), t.TPS(100), t.TPS(200), t.TPS(600))
}

// Close ...
func (t *TPSTracker) Close() error {
	t.once.Do(func() {
		close(t.closing)
		t.wg.Wait()
		t.w = nil
	})
	return nil
}
