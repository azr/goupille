// Goupille is the french for a grenade's pin.
//
// Pull the pin, throw stuff and run.
//
// This package it is intended for you to make
// sure every goroutine left after anyone pulled the pin, simply.
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
	"sync"
)

type Goupille interface {

	// Add a worker
	Add()

	// Call Pull to tell everyone something went wrong
	Pull(reason error)

	// if this is triggered, it means you should leave if possible
	Tick() chan struct{}

	// Done working
	Done()

	// a true hero waits for everyone before leaving
	Wait() error
}

//My version of la goupille
type Pin struct {
	m      sync.Mutex
	g      sync.WaitGroup
	dying  chan struct{}
	reason error
}

// New safety Pin !
func New() *Pin {
	return &Pin{
		m:      sync.Mutex{},
		g:      sync.WaitGroup{},
		dying:  nil,
		reason: nil,
	}
}

// Add a worker
func (g *Pin) Add() {
	g.g.Add(1)
}

// Worker finished
func (g *Pin) Done() {
	g.g.Done()
}

// Pull the pin.
//
// Hopefully this tells everyone attached to start leaving
// before explosion...
func (g *Pin) Pull(reason error) {
	g.m.Lock()
	defer g.m.Unlock()

	if g.reason == nil {
		g.reason = reason
	}

	select {
	//a nil chan never gets triggered
	//but a closed chan always does
	case <-g.dying:
	default:
		g.dying = make(chan struct{})
		close(g.dying)
	}
}

// Tick tac motherlover !
// Let's stop working and start leaving.
func (g *Pin) Tick() chan struct{} {
	return g.dying
}

// Wait until everyone leaves
// and gives the termination reason.
func (g *Pin) Wait() error {
	g.g.Wait()
	g.m.Lock()
	defer g.m.Unlock()
	return g.reason
}
