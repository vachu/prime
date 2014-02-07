package main

import "os"
import "fmt"
import "flag"

import "primelib"

var print = flag.Uint("print", 0, "prints first 'n' primes")
var test = flag.Bool("test", false, "tests the given number/s (in <STDIN>) for primality")

func main() {
	flag.Parse()

	if *print > 0 {
		fmt.Fprintf(os.Stderr, "INFO: Printing %d primes...\n", *print)

		cnt := primelib.WritePrimes(os.Stdout, uint32(*print), "\n")
		if cnt < uint32(*print) {
			fmt.Fprintf(os.Stderr, "WARNING: printed %d primes\n", cnt)
		}
	}
	if !(*test) {
		return
	}

	fmt.Fprintf(os.Stderr, "INFO: input number for Primality Test; ")
	fmt.Fprintf(os.Stderr, "<Ctrl-D> to quit...\n")
	testPrime()
}

func testPrime() {
	for number := uint64(0); getNumber(&number); {
		fmt.Print(number, ": ")

		switch primeFactor := primelib.GetFirstPrimeFactor(number); {
			case primeFactor == 0:
				fmt.Println("cannot test reliably.")

			case primeFactor == number:
				fmt.Println("prime")

			default:
				fmt.Print("composite - divisible by ", primeFactor)
				for number != primeFactor {
					lastPrimeFactor := primeFactor
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
