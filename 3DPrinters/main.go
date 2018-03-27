package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//ThreeDPrinter ...
type ThreeDPrinter interface {
	Print()
}

//HackerPrinter ...
type HackerPrinter struct {
	name      string
	tired     <-chan time.Time
	completed chan struct{}
	timesUsed int
}

func main() {
	//Create a pool of 3 3D printers (each 3d printers is a goroutine)
	// They are waiting in a channel
	pool := NewPool(2)

	//Create the 7 hackers and use a waitgroup to know when they finish printing
	//or they leave
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		hp := HackerPrinter{name: fmt.Sprint(i), tired: time.After(5 * time.Second), completed: make(chan struct{})}
		go func() {
			enqueue(pool, &hp)
			defer wg.Done()
		}()
	}

	wg.Wait() //wait for the 7 hackers to print or leave
	fmt.Println("Shutting down")
	//All hackers started already their printing jobs (or already left), we can shutdown the pool of printers
	pool.Shutdown()
}

//enqueue Puts a hacker to wait into the queue until a printer is ready
//The hacker will leave the queue if he is waiting more than specified timeout before
//a printer is ready
func enqueue(pool *Pool, hp *HackerPrinter) {
	hp.timesUsed++
	fmt.Printf("Hacker %s enters the queue for %d time\n", hp.name, hp.timesUsed)
	//If timeout happens before the send in queue, then "the hacker leaves"
	select {
	case pool.GetQueue() <- hp:
		hp.waitForFinishPrint(pool)
		return
	case <-hp.tired:
		fmt.Printf("Hacker %s leaves the hackathon tired of waiting\n", hp.name)
		return
	}
}

//Print ...
func (hp *HackerPrinter) Print() {
	fmt.Printf("Hacker %s starts printing...\n", hp.name)
	max := 10
	min := 1
	rnd := rand.Intn(max-min) + min
	d := time.Duration(rnd)
	time.Sleep(d * time.Second)
	fmt.Printf("Hacker %s finish printing\n", hp.name)
	hp.completed <- struct{}{}
}

//waitForFinishPrint waits until print job is completed and then
//if it's the first time goes to the queue again
func (hp *HackerPrinter) waitForFinishPrint(pool *Pool) {
	<-hp.completed
	if hp.timesUsed < 2 {
		hp.tired = time.After(5 * time.Second)
		enqueue(pool, hp)
	} else {
		return
	}
}
