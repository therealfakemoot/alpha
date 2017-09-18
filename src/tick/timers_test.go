package discord

import (
	"fmt"
	"testing"
	"time"
)

func TestBasicTicker(t *testing.T) {
	// t.Skip("TickOnce not implemented")

	var counter = 0
	var i = 0

	f := func(to *Ticker) {

		fmt.Printf("Iteration: %d\n", i)
		if i >= 5 {
			to.Done()
		}

		counter++
		i++
	}

	cleanUp := func(to *Ticker) {
		return
	}

	ticker := NewTicker(1*time.Second, f, cleanUp)

	time.AfterFunc(6*time.Second, func() {
		ticker.Done()
	})

	select {
	case q := <-ticker.Quit:
		if counter != 5 && q {
			t.Errorf("\nIncorrect `counter` value.\nExpected: %d | Received: %d", 5, counter)
		}
	}
}

func TestTickerStop(t *testing.T) {
	t.Skip("TickerStop not impelemnted.")
}
