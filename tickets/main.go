package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	tickets := 500
	ticketChan := make(chan int)
	doneChan := make(chan bool)

	go manageTicket(ticketChan, doneChan, &tickets)

	for userId := range 2000 {
		wg.Add(1)
		go buyTicket(&wg, ticketChan, userId)
	}

	wg.Wait()
	doneChan <- true
}

func buyTicket(wg *sync.WaitGroup, ticketChan chan int, userId int) {
	defer wg.Done()
	ticketChan <- userId
}

func manageTicket(ticketChan chan int, doneChan chan bool, tickets *int) {
	for {
		select {
		case user := <-ticketChan:
			if *tickets > 0 {
				*tickets--
				fmt.Printf("Ticket purchased by user %d. Tickets remaining: %d\n", user, *tickets)
			} else {
				fmt.Printf("USer %d found no tickets\n", user)
			}
		case <-doneChan:
			fmt.Printf("Tickets remaining: %d\n", *tickets)
		}
	}
}
