package main

import "os"
import "fmt"
import "flag"

import "primelib"

var generate = flag.Uint("gen", 0, "generates first 'n' primes")
var test = flag.Bool("test", false, "tests the given number/s (in <STDIN>) for primality")
var inputFile = flag.String("file", "", "loads the first 'n' primes from the specified file - 'n' is specified with '-gen' flag")
var noPrint = flag.Bool("noprint", false, "suppress printing primes")

func main() {
	flag.Parse()

	if inputFile != nil && len(*inputFile) > 0 {
		cnt, err := primelib.LoadPrimes(*generate, *inputFile)
		if cnt < *generate {
			fmt.Fprint(os.Stderr, "WARNING: error encountered after ")
			fmt.Fprintln(os.Stderr, "loading", cnt, "prime/s -", err)
		}
	}
	primelib.GeneratePrimes(*generate)
	if primelib.GetPrimeCount() == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: No primes generated")
		return
	}

	if !(*test) {
		if !(*noPrint) {
			printPrimes()
		} else {
			largestPrime := primelib.GetMaxPrime()
			fmt.Println("Largest prime generated: ", largestPrime)
		}
		return
	}

	var number uint64
	maxPrime := primelib.GetMaxPrime()
	threshold := maxPrime * maxPrime
	for getNumber(&number) {
		fmt.Print(number, ": ")

		if number > threshold {
			fmt.Print("cannot test reliably.  ")
			fmt.Println("Maximum testable number =", threshold)
			continue
		}

		primeFactor := primelib.GetFirstPrimeFactor(number)
		if primeFactor == number {
			fmt.Println("prime")
		} else {
			fmt.Print("composite - divisible by ", primeFactor)
			lastPrimeFactor := primeFactor

			for ; number != primeFactor; lastPrimeFactor = primeFactor {
				number /= primeFactor
				primeFactor = primelib.GetFirstPrimeFactor(number)
				if primeFactor != lastPrimeFactor {
					fmt.Print(", ", primeFactor)
				}
			}
			fmt.Println("")
		}
	}
}

func getNumber(number *uint64) bool {
	_, err := fmt.Scanf("%d", number)

	return err == nil
}

func printPrimes() {
	var i uint
	for p := primelib.GetPrime(i); p != 0; p = primelib.GetPrime(i) {
		fmt.Println(p)
		i++
	}
}
