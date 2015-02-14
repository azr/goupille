package goupille

import (
	"fmt"
	"os"
	"os/signal"
)

//Calls os.Notify which will trigger ending (Pull) on selected `os.Signal`(s).
//If you pass nil, any signal will trigger ending.
//
//This starts a goroutine
func (g *Pin) Notify(sig ...os.Signal) {
	c := make(chan os.Signal)
	signal.Notify(c, sig...)
	go func() {
		for sig := range c {
			g.Pull(fmt.Errorf("Recieved signal %s", sig))
		}
	}()
}
