package regulator

import (
	"log"
	"time"
)

type Regulator struct {
	In      chan interface{}
	Out     chan interface{}
	ticksPs int
	stop    chan struct{}
}

// New return a new regulator with given chan sizes and ticks per second
func New(inSize, outSize, ticksPs int) (r *Regulator) {
	r = &Regulator{
		In:      make(chan interface{}, inSize),
		Out:     make(chan interface{}, outSize),
		stop:    make(chan struct{}),
		ticksPs: ticksPs,
	}
	return
}

// Start copies elements from in channel to the out channel
func (r *Regulator) Start(opsPerSecond int) {
	r.run(opsPerSecond)
}

// Stop copying
func (r *Regulator) Stop() {
	r.stop <- struct{}{}
}

func (r *Regulator) run(opsPerSecond int) {
	var ticks int
	if opsPerSecond < r.ticksPs {
		ticks = opsPerSecond
	} else {
		ticks = r.ticksPs
	}

	interval := time.Second / time.Duration(ticks)
	elems := opsPerSecond / ticks
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			go r.copier(elems, interval)
		case <-r.stop:
			ticker.Stop()
			// Clean any pending ticker message
			select {
			case <-ticker.C:
			default:
			}
			return
		}
	}
}

func (r *Regulator) copier(n int, interval time.Duration) {
	for i := 0; i < n; i++ {
		select {
		case e := <-r.In:
			select {
			case r.Out <- e:
			default:
				log.Println("Overflow in output channel")
			}
		default:
			log.Println("No elements in input channel")
		}
	}
}
