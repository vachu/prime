package main

import (
	"fmt"
	"primelib/v4"
)

func main() {
	_, out, diag := primelib.ListPrimesBetween(0, 10000)

	if out != nil {
		for output := range out {
			fmt.Println(output)
		}
	}
	fmt.Println("DEBUG: printed", <-diag, "primes")
}

func drain(ch chan interface{}) {
	if ch != nil {
		for data := range ch {
			_ = data
		}
	}
}
