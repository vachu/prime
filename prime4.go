package main

import (
	"fmt"
	"os"
	"primelib/v4"
)

func main() {
	//cmdLine := "?"
	//cmdLine := "ListPrimesBetween 0 20"
	cmdLine := "ListPrimes 100"
	in, out, diag := primelib.DoCmd(cmdLine)

	ctr, breakout := 0, 0
	if out != nil {
		for output := range out {
			fmt.Println(output)
			ctr++
			if ctr == breakout {
				fmt.Println("DEBUG: encountered BREAKOUT -- ctr =", ctr)
				close(in)
				break
			}
		}
		fmt.Println("DEBUG: ctr =", ctr)
		<-drain(out)
	}
	fmt.Fprintln(os.Stderr, <-diag)
}

func drain(ch chan interface{}) (done chan interface{}) {
	done = make(chan interface{})
	go func() {
		defer close(done)

		ctr := 0
		if ch != nil {
			for data := range ch {
				_ = data
				ctr++
			}
		}
		fmt.Println("DEBUG: from go func() - drain count =", ctr)
	}()
	return
}
