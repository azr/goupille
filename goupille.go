// Goupille is the french for a grenade's pin.
//
// Pull the pin, throw stuff and run.
//
// This package it is intended for you to make
// sure every goroutine left after anyone pulled the pin.
//
// It is pretty convenient when you have multiple goroutines
// doing things and one of them could fail. So you
// can make sure everyone says goodbye before leaving, and then shut down properly.
//
// Goupille is meant to be used as an interface, you can have your own implementation of it and use it everywhere.
//
// For example, my implementation (Pin) starts handling system os.Signal with Notify.
package goupille

import (
	"os"
	"os/signal"
	"sync"
)

type Goupille interface {

	//1 Attatch a goroutine to it
	Attach()

	//2 Something went wrong, let's pull the goupille
	//  ( you should also throw the grenade)
	Pull(reason Stringer)

	//3 tell everyone the chemical reaction started
	Tick() chan struct{}

	//4 yay that goroutine left
	Detach()

	//5 a true hero waits for everyone before leaving
	Wait() Stringer
}

//Stringer is already definer in fmt !
//but a real Goupille doesn't need fmt,
//question de fierté.
type Stringer interface {
	String() string
}

//My version of la goupille
type Pin struct {
	m       sync.Mutex
	g       sync.WaitGroup
	waiting chan struct{}
	reason  Stringer
}

// New safety Pin !
func New() *Pin {
	return &Pin{
		m:       sync.Mutex{},
		g:       sync.WaitGroup{},
		waiting: nil,
		reason:  nil,
	}
}

// Attach a string to the Pin.
// Remember, if you pull it, start leaving
func (g *Pin) Attach() {
	g.g.Add(1)
}

// Cuts the string (please use legs to run away)
func (g *Pin) Detach() {
	g.g.Done()
}

// Pull the pin.
//
// Hopefully this tells everyone attached to start leaving
// before explosion...
func (g *Pin) Pull(reason Stringer) {
	g.m.Lock()
	defer g.m.Unlock()

	if g.reason == nil {
		g.reason = reason
	}

	select {
	//a nil chan never gets triggered
	//but a closed chan always does
	case <-g.waiting:
	default:
		g.waiting = make(chan struct{})
		close(g.waiting)
	}
}

// Tick tac motherlover !
// Leave !!!
func (g *Pin) Tick() chan struct{} {
	return g.waiting
}

//Calls os.Notify which will trigger ending (Pull) on selected `os.Signal`(s).
//If you pass nil, any signal will trigger ending.
//
//This starts a goroutine
func (g *Pin) Notify(sig ...os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, sig...)
	go func() {
		for sig := range c {
			g.Pull(sig)
		}
	}()
}

// Wait until everyone leaves
// and gives the termination reason.
func (g *Pin) Wait() Stringer {
	g.g.Wait()
	g.m.Lock()
	defer g.m.Unlock()
	return g.reason
}
