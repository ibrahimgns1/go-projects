package main

import (
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

type AutoClicker struct {
	isRunning bool
	interval  time.Duration
	button    string
	quitCh    chan struct{}
	mu        sync.Mutex
}

func NewAutoClicker() *AutoClicker {
	return &AutoClicker{
		isRunning: false,
		interval:  time.Millisecond * 100,
		button:    "left",
		quitCh:    make(chan struct{}),
	}
}

func (ac *AutoClicker) Start() {
	ac.mu.Lock()
	if ac.isRunning {
		ac.mu.Unlock()
		return
	}
	ac.isRunning = true
	ac.mu.Unlock()

	go func() {

		ac.mu.Lock()
		currentInterval := ac.interval
		currentButton := ac.button
		ac.mu.Unlock()

		ticker := time.NewTicker(currentInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ac.quitCh:
				return
			case <-ticker.C:
				robotgo.Click(currentButton)
			}
		}
	}()
}

func (ac *AutoClicker) Stop() {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	if !ac.isRunning {
		return
	}

	ac.quitCh <- struct{}{}
	ac.isRunning = false
}

func (ac *AutoClicker) SetInterval(ms int) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.interval = time.Duration(ms) * time.Millisecond
}

func (ac *AutoClicker) SetButton(btn string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.button = btn
}

func (ac *AutoClicker) IsRunning() bool {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	return ac.isRunning
}
