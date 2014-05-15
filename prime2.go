package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

import "primelib/v3"

const countFlagUsage = `x:y
	Prints the count of all primes between 'x' and 'y' where 'x' <= 'y'
`

const printFlagUsage = `{n | x:y}
    Prints first 'n' primes or all primes between 'x' and 'y' where 'x' <= 'y'
`
const testFlagUsage = `
    Tests the given number/s (in <STDIN>) for primality
`

var count = flag.String("count", "", countFlagUsage)
var print = flag.String("print", "", printFlagUsage)
var test = flag.Bool("test", false, testFlagUsage)

func main() {
	flag.Parse()
	if len(*count) > 0 {
		countPrime()
	} else if len(*print) > 0 {
		printPrime()
	}
	if *test {
		fmt.Fprintf(os.Stderr, "INFO: input number for Primality Test; ")
		fmt.Fprintf(os.Stderr, "<Ctrl-D> to quit...\n")
		testPrime()
	}
}

func countPrime() {
	_, from, to := parsePrintFlag(*count)
	ch := make(chan uint32, 1000)
	switch {
	case from <= to && to > 0:
		go primelib.ListPrimesBetween(ch, from, to)
		cnt := uint32(0)
		for _ = range ch {
			cnt++
		}
		fmt.Println("\n", cnt)

	default:
		fmt.Fprintf(os.Stderr, "ERROR: Invalid 'print' flag value(s)\n")
	}
}

func printPrime() {
	cnt, from, to := parsePrintFlag(*print)
	ch := make(chan uint32, 1000)
	switch {
	case from <= to && to > 0:
		go primelib.ListPrimesBetween(ch, from, to)
		for p := range ch {
			fmt.Println(p)
		}

	case cnt > 0:
		go primelib.ListPrimes(ch, cnt)
		for p := range ch {
			fmt.Println(p)
		}

	default:
		fmt.Fprintf(os.Stderr, "ERROR: Invalid 'print' flag value(s)\n")
	}
}

func parsePrintFlag(flagValue string) (cnt, from, to uint32) {
	cnt = 0
	from = 0
	to = 0
	r1 := regexp.MustCompile(`^\d+$`)
	r2 := regexp.MustCompile(`^\d+:\d+`)
	r3 := regexp.MustCompile(`:`)
	if r1.MatchString(flagValue) {
		c, err := strconv.ParseUint(flagValue, 0, 32)
		if err == nil {
			cnt = uint32(c)
		}
	} else if r2.MatchString(flagValue) {
		rangeNums := r3.Split(flagValue, -1)
		f, e1 := strconv.ParseUint(rangeNums[0], 0, 32)
		t, e2 := strconv.ParseUint(rangeNums[1], 0, 32)
		if e1 == nil {
			from = uint32(f)
		}
		if e2 == nil {
			to = uint32(t)
		}
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
	_, err := fmt.Scanln(number)

	return err == nil
}
