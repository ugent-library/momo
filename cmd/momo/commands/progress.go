package commands

import (
	"fmt"
	"time"
)

type progress struct {
	start time.Time
	every int
	count int
	i     int
}

func newProgress(every int) *progress {
	p := &progress{start: time.Now(), every: every}
	return p
}

func (p *progress) rate() int {
	dur := time.Since(p.start)
	return int(1000.0 / float64(dur/time.Millisecond) * float64(p.count))
}

func (p *progress) inc() {
	p.count++
	p.i++
	if p.i == p.every {
		p.i = 0
		fmt.Printf("\r%d", p.count)
	}
}

func (p *progress) done() {
	fmt.Printf("\r%d done (%d/sec)\n", p.count, p.rate())
}
