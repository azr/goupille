package goupille_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/azr/goupille"
)

//Have multiple workers do stuff
func TestExampleWorkers(t *testing.T) {
	g := goupille.New()

	f := func() {
		g.Add()
		g.Done()
		for {
			select {
			case <-time.After(10 * time.Nanosecond):
				//Doing stuff//
			case <-g.Tick():
				fmt.Println("Tick !")
				g.Pull(nil) // pull for no reason
				return
			}
		}
	}

	go f()
	go f()
	go f()

	time.Sleep(time.Nanosecond)
	g.Pull(nil)
	g.Wait()
}

func TestBusyWorkers(t *testing.T) {
	g := goupille.New()
	length := 100
	c := make(chan int, length)

	go func() {
		g.Add()
		defer g.Done()
		for i := 0; ; i++ {
			select {
			case <-time.After(time.Nanosecond):
				c <- 2
				if i == length-1 {
					g.Pull(fmt.Errorf("handlable error"))
					g.Pull(fmt.Errorf("Not this error"))
				}
			case <-g.Tick():
				close(c) //stopped producing
				return
			}
		}
	}()

	done := 0
	go func() {
		g.Add()
		defer g.Done()
		for {
			select {
			case _, hasMore := <-c:
				if !hasMore {
					return //No more stuff to process
				}
				done++
			}
		}
	}()

	time.Sleep(time.Nanosecond)

	err := g.Wait()
	fmt.Printf("Ended after '%s'\n", err)
	if done != length {
		t.Error("It didn't do enough !")
	}
}

func asserter(_ goupille.Goupille) {
}

//lets prevent builds if goupille is not a real Goupille
func goupilleIsNotAGoupille() {
	asserter(goupille.New())
}
