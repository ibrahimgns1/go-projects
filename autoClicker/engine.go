package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

type AutoClicker struct {
	isRunning   bool
	interval    time.Duration
	button      string
	useLocation bool
	locationX   int
	locationY   int
	quitCh      chan struct{}
	mu          sync.Mutex
}

func NewAutoClicker() *AutoClicker {
	return &AutoClicker{
		isRunning:   false,
		interval:    time.Millisecond * 100,
		button:      "left",
		useLocation: false,
		locationX:   0,
		locationY:   0,
		quitCh:      make(chan struct{}),
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
		useLocation := ac.useLocation
		locationX := ac.locationX
		locationY := ac.locationY
		ac.mu.Unlock()

		ticker := time.NewTicker(currentInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ac.quitCh:
				return
			case <-ticker.C:
				if useLocation {
					fmt.Printf("Clicking at %d, %d with %s button\n", locationX, locationY, currentButton)
					robotgo.Move(locationX, locationY)
					robotgo.Click(currentButton)
				} else {
					robotgo.Click(currentButton)
				}
			}
		}
	}()
}

func (ac *AutoClicker) Stop() {
	ac.mu.Lock()
	if !ac.isRunning {
		ac.mu.Unlock()
		return
	}
	ac.isRunning = false
	ac.mu.Unlock()

	select {
	case ac.quitCh <- struct{}{}:
	default:
	}
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

func (ac *AutoClicker) SetUseLocation(use bool) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.useLocation = use
}

func (ac *AutoClicker) SetLocation(x, y int) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.locationX = x
	ac.locationY = y
}
