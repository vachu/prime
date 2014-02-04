package main

import "primelib"
import "os"

func main() {
	primelib.WritePrimes(os.Stdout, 20, " ****\n");
}
