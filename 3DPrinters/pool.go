package main

import (
	"sync"
)

//Pool ...
type Pool struct {
	printers chan ThreeDPrinter
	wg       sync.WaitGroup
}

//NewPool initializes a new instace of Pool
func NewPool(numPrinters int) *Pool {
	p := Pool{
		printers: make(chan ThreeDPrinter),
	}

	p.wg.Add(numPrinters)

	for i := 0; i < numPrinters; i++ {
		go func() {
			for printer := range p.printers {
				printer.Print()
			}
			p.wg.Done()
		}()
	}

	return &p
}

//GetQueue get channel to enque printers
func (p *Pool) GetQueue() chan ThreeDPrinter {
	return p.printers
}

//Shutdown ...
func (p *Pool) Shutdown() {
	//no more sends into the channel are accepted
	close(p.printers)
	//wait until all printer goroutines are finished
	p.wg.Wait()
}
