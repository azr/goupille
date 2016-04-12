Deprecated
==========
Warning ! I made this pkg before I knew of the context pkg.
Even if goupille work, the net/context pkg is better solution for this.

Goupille
========

[![GoDoc](https://godoc.org/github.com/azr/goupille?status.png)](https://godoc.org/github.com/azr/goupille)
[![Build Status](https://travis-ci.org/azr/goupille.svg?branch=master)](https://travis-ci.org/azr/goupille)

This package it is intended for you to make
sure every goroutine left after anyone pulled the pin.

**Goupille** is the french for a grenade's pin.

Pull the pin, throw stuff and run.


It is pretty convenient when you have multiple goroutines
doing things and one of them could fail. So you
can make sure everyone says goodbye before leaving, and then shut down properly.

Goupille is meant to be used as an interface, you can have your own implementation of it and use it everywhere.

For example, my implementation (Pin) starts handling system os.Signal with Notify.
