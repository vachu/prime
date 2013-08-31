package main

import "fmt"
import "time"

const MAX_GO_ROUTINES = 2 

func main() {
	//channelTest1()
	cancellationChannel()
}

func goProc(id int, cancelChannel chan int) {
	for ctr := 1; ; ctr++ {
		select {
			case _ = <-cancelChannel:
				fmt.Println("From goProc #", id, ": Cancelled")
				return
			default:
				fmt.Println("From goProc #", id, ": ctr =", ctr)
				time.Sleep(time.Second)
		}
	}
}

func cancellationChannel() {
	cancelChannel := make(chan int, MAX_GO_ROUTINES)
	for i := 0; i < MAX_GO_ROUTINES; i++ {
		go goProc(i + 1, cancelChannel)
	}

	time.Sleep(time.Second * 10)
	for i := 0; i < MAX_GO_ROUTINES; i++ {
		cancelChannel <- 1
	}
	fmt.Println("Cancellation Token sent")
	time.Sleep(time.Second * 1)
}

func delayedProc(output chan int, data int) {
	fmt.Println("From Goroutine #", data)
	output <- data
}

func channelTest1() {
	output := make(chan int)
	for i := 0; i < MAX_GO_ROUTINES; i++ {
		go delayedProc(output, i+1)
	}

	fmt.Println("Waiting for all goroutines to finish...")
	for i := 0; i < MAX_GO_ROUTINES; i++ {
		data := <- output
		fmt.Println("Goroutine #", data, "completed")
	}
}
