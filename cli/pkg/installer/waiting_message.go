package installer

import (
	"log"
	"time"
)

type WaitingMessage struct {
	exit chan<- int
}

func NewWaitingMessage(name string, delay time.Duration) *WaitingMessage {
	exit := make(chan int, 1)
	go func() {
		iterations := 0
		for {
			select {
			case <-exit:
				return
			case <-time.After(delay):
				iterations++
				log.Printf("Waiting for %s (%v elapsed)", name, delay*time.Duration(iterations))
			}
		}
	}()
	return &WaitingMessage{
		exit: exit,
	}
}

func (w *WaitingMessage) Stop() {
	defer close(w.exit)
	w.exit <- 0
}
