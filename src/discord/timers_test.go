package discord

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBasicTicker(t *testing.T) {
	// t.Skip("TickOnce not implemented")

	var counter = 1
	var i = 0

	mu := sync.Mutex{}

	f := func(to *Ticker) {
		fmt.Println("Iteration: " + string(i))
		mu.Lock()
		defer mu.Unlock()

		if i >= 5 {
			to.Done()
		}

		counter++
		i++
	}

	cleanUp := func(to *Ticker) {
		return
	}

	ticker := NewTicker(1, f, cleanUp)

	time.AfterFunc(6, func() {
		ticker.Done()
	})

	select {
	case q := <-ticker.Quit:
		if counter != 5 && q {
			t.Errorf("Incorrect `counter` value.")
			t.Errorf("Expected: %d | Received: %d", 5, counter)
		}
	}
}

func TestTickerStop(t *testing.T) {
	t.Skip("TickerStop not impelemnted.")
}
