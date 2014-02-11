package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)
import "primelib"

const printFlagUsage = `{n | x:y}
    Prints first 'n' primes or all primes between 'x' and 'y' where 'x' <= 'y'
`
const testFlagUsage = `
    Tests the given number/s (in <STDIN>) for primality
`

var print = flag.String("print", "6542", printFlagUsage)
var test = flag.Bool("test", false, testFlagUsage)

func main() {
	flag.Parse()

	if len(*print) > 0 {
		printPrime()
	}
	if !(*test) {
		return
	}

	fmt.Fprintf(os.Stderr, "INFO: input number for Primality Test; ")
	fmt.Fprintf(os.Stderr, "<Ctrl-D> to quit...\n")
	testPrime()
}

func printPrime() {
	cnt, from, to := parsePrintFlag(*print)
	fmt.Fprintf(os.Stderr, "DEBUG: cnt = %d, from = %d, to = %d\n", cnt, from, to)
	switch {
	case cnt > 0:
		fmt.Fprintf(os.Stderr, "INFO: Printing %d primes...\n", cnt)
		ctr := primelib.WritePrimes(os.Stdout, cnt, "\n")
		if ctr < uint32(cnt) {
			fmt.Fprintf(os.Stderr, "WARNING: printed %d primes\n", ctr)
		}

	case from < to:
		fmt.Fprintf(os.Stderr, "INFO: Printing primes between %d and %d...\n", from, to)
		ctr := primelib.WritePrimesBetween(os.Stdout, from, to, "\n")
		fmt.Fprintf(os.Stderr, "INFO: printed %d primes\n", ctr)

	default:
		fmt.Fprintf(os.Stderr, "Invalid 'print' flag value(s)")
	}
}

func parsePrintFlag(flagValue string) (cnt, from, to uint32) {
	cnt = 0
	from = 0
	to = 0
	r1 := regexp.MustCompile(`^\d+$`)
	r2 := regexp.MustCompile(`^\d+:\d+$`)
	r3 := regexp.MustCompile(`:`)
	if r1.MatchString(flagValue) {
		c, err := strconv.ParseUint(flagValue, 0, 32)
		if err == nil {
			cnt = uint32(c)
		}
	}
	if r2.MatchString(flagValue) {
		rangeNums := r3.Split(flagValue, -1)
		f, e1 := strconv.ParseUint(rangeNums[0], 0, 32)
		t, e2 := strconv.ParseUint(rangeNums[1], 0, 32)
		if e1 == nil {
			from = uint32(f)
		}
		if e2 == nil {
			to = uint32(t)
		}
		fmt.Fprintf(os.Stderr, "DEBUG: from = %d, to = %d\n", from, to)
	}
	return
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
