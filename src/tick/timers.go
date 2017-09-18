package discord

import (
	"time"
)

// Ticker describes a value capable of performing one task at regular intervals with support for pausing/resuming with Fire()/Pause() and finalizing via Done()
//
// Fire() begins ( or resumes )
// returns no value, it simply executes a side effect. Users must mind their state and use mutxes when appropriate.
type Ticker struct {
	F       TimerFunc
	CleanUp CleanupFunc
	Dur     time.Duration
	Timer   *time.Timer
	Quit    chan bool
}

// TimerFunc is a named type describing a callback that can interact with its parent timer.
type TimerFunc func(*Ticker)

// CleanupFunc is similar to TimerFunc, but is only called once, when Done() is called.
type CleanupFunc func(*Ticker)

func (to *Ticker) refresh() {
	to.Timer = NewTicker(to.Dur, to.F, to.CleanUp).Timer
}

// Fire executes TimerFunc in a goroutine and defers the refresh() method to start the next cycle.
func (to *Ticker) Fire() {
	go to.F(to)
	defer to.refresh()
}

// Stop terminates the timer and invokes Done() and the CleanupFunc.
func (to *Ticker) Stop() {
	to.Timer.Stop()
	to.Done()
}

// Done executes the CleanupFunc and sends a true value on the Quit channel.
func (to *Ticker) Done() {
	to.CleanUp(to)
	to.Quit <- true
}

// NewTicker returns a new timed discord event.
func NewTicker(d time.Duration, f TimerFunc, c CleanupFunc) *Ticker {
	ticker := &Ticker{}

	doFunc := func() {
		ticker.F(ticker)
		defer ticker.refresh()
	}
	t := time.AfterFunc(d, doFunc)

	ticker.Dur = d
	ticker.F = f
	ticker.Timer = t
	ticker.CleanUp = c
	ticker.Quit = make(chan bool, 1)
	return ticker
}
