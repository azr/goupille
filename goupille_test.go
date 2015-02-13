package goupille_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/azr/goupille"
)

//Have multiple workers do stuff
//
func TestExampleWorkers(t *testing.T) {
	g := goupille.New()

	fmt.Println("Kikoo")
	f := func() {
		g.Attach()
		for {
			select {
			case <-time.After(100 * time.Microsecond):
				fmt.Println("Doing stuff")
			case <-g.Tick():
				fmt.Println("Tick !")
				g.Detach()
				g.Pull(nil) // pull for no reason
				return
			}
		}
	}

	go f()
	go f()
	go f()

	time.Sleep(time.Microsecond * 200)
	g.Pull(nil)
	g.Wait()
	fmt.Println("Gone !!")
}

func asserter(_ goupille.Goupille) {
}

//lets prevent builds if goupille is not a real Goupille
func goupilleIsNotAGoupille() {
	asserter(goupille.New())
}
