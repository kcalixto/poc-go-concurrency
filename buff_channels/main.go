package main

import (
	"fmt"
	"time"
)

func main() {
	orders := make(chan string, 2) // buffered chan

	// customer placing orders
	go func() {
		for i := range 5 {
			order := fmt.Sprintf("Coffee order #%d", i)
			orders <- order // only blocks when buffer is full
			fmt.Println("Placed:", order)
		}
		close(orders)
	}()

	// Barista processing orders
	for order := range orders { // blocks when buffer is empty
		// fmt.Println("Processing: ", order)
		time.Sleep(2 * time.Second)
		fmt.Printf("☕️ served: %s\n", order)
	}
}
