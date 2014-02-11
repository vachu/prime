package main

import (
	"fmt"
	"os"
	"primelib"
)

func main() {
	const req_cnt = 20000

	cnt := primelib.WritePrimes(os.Stdout, req_cnt, "\n")
	if cnt != req_cnt {
		fmt.Fprintf(os.Stderr, "WARNING: got only %d primes\n", cnt)
	}
}
