/*
A much simplified package that provides APIs for listing and testing primes.
Internally generates first 1 million unsigned primes that are used for all the
operations, at startup / init.

NOTE: Number '1' is not considered a prime by this package.  Negative Numbers /
Primes are not considered either.

All the public / exported methods in this package that return the channel triad
(in, out, diag) run their main logic concurrently while returning the aforesaid
triad immediately.  The expected usage of each of the channels in the triad is
as follows:

- in  : the caller could supply input data for processing and / or abort the
spawned concurrent thread by just closing this channel

- out : the caller could get the complete output by reading this channel until
it is closed by the Method Implementation

- diag: the caller could get diagnostic output (if any) through this channel
until it is closed by the Method Implementation

This Design Pattern is tentatively named as "Std3io" pattern since it is very
much analogous to the (STDIN, STDOUT, STDERR) triad which is well known in the
console / terminal applications

IMPORTANT NOTE: Method implementations following this Std3io pattern can return
nil for 'in' and 'diag' channels; that only means that the implementation is not
providing any functionality based on those channels.  If the 'out' channel is
nil, then the caller should presume that something is wrong with the method call.
Anyways, the relevant Method Documentation should clearly specify the gory
details.
*/
package primelib

import (
	"fmt"
	"math"
)

const maxPrimeCount = 10 * 1000

var arrPrimes [maxPrimeCount]uint32

func init() {
	arrPrimes[0] = 2
	primeCount := 1
	for number := uint32(3); primeCount < maxPrimeCount; number += 2 {
		if isPrime(uint64(number)) {
			arrPrimes[primeCount] = number
			primeCount++
		}
	}
}

func isPrime(number uint64) bool {
	isComposite := false
	stopLimit := uint64(math.Trunc(math.Sqrt(float64(number))))
	for i := 0; !isComposite && uint64(arrPrimes[i]) <= stopLimit; i++ {
		isComposite = (number%uint64(arrPrimes[i]) == 0)
	}
	return !isComposite
}

func makeChanTrio(writeBuffSize uint32) (in, out, diag chan interface{}) {
	in = make(chan interface{}, 1)
	out = make(chan interface{}, writeBuffSize)
	diag = make(chan interface{}, writeBuffSize)
	return
}

// Lists out the first 'cnt' primes onto the 'out' channel.
//
// This method conforms to the Std3io pattern.  It returns valid 'in', 'out'
// and 'diag' channels.  The main output of the list of first 'cnt' primes
// is written onto the 'out' channel and the count of which is written onto
// the 'diag' channel.  Each of the primes and the final count are of type
// 'uint32'
//
// The caller can abort the concurrent execution by closing the 'in' channel
func ListPrimes(cnt uint32) (in, out, diag chan interface{}) {
	in, out, diag = makeChanTrio(10 * 1000) // new diag channel created
	go listPrimes(cnt, in, out, diag)
	return
}

func listPrimes(cnt uint32, in, out, diag chan interface{}) {
	defer close(out)
	defer close(diag)

	status := "OK"
	primePrintCnt := uint32(0)
	number := uint64(arrPrimes[maxPrimeCount-1])
	stopLimit := number * number
	if number%2 == 0 {
		number -= 1
	}
MainLoop:
	for i := uint32(0); i < cnt && number < stopLimit; i++ {
		select {
		case _, isOpen := <-in:
			if !isOpen {
				status = "ABORTED"
				break MainLoop
			}
		default:
			if i < uint32(maxPrimeCount) {
				out <- arrPrimes[i]
				primePrintCnt++
			} else {
				for number += 2; number < stopLimit; number += 2 {
					if isPrime(number) {
						out <- number
						primePrintCnt++
						break
					}
				}
			}
		} // select
	} // for ...
	diag <- fmt.Sprintf("%s: Listed %d primes", status, primePrintCnt)
}

// List all 32-bit primes between >= 'from' AND <= 'to' onto the returned
// 'out' channel.  The total count of the primes written to 'out' is written
// onto 'diag' channel.  The caller can abort the concurrent execution by
// closing the 'in' channel
//
// This method conforms to the Std3io pattern. The channels 'in', 'out' and
// 'diag' are immediately returned to the caller. In case of an error (mostly
// illegal 'from' and 'to' values) then only the 'diag' channel is valid while
// the rest are 'nil'
//
// Each prime written onto the 'out' channel is of type uint64 and the count
// written onto 'diag' is of type uint32
func ListPrimesBetween(from, to uint32) (in, out, diag chan interface{}) {
	diag = make(chan interface{}, 1)
	defer close(diag)

	if to < 2 || from > to {
		diag <- fmt.Sprintf("ERROR: Illegal from (%d) / to (%d) values", from, to)
		return
	}

	in, out, diag = makeChanTrio(10 * 1000)
	go listPrimesBetween(from, to, in, out, diag)
	return
}

func listPrimesBetween(from, to uint32, in, out, diag chan interface{}) {
	defer close(out)
	defer close(diag)

	status := "OK"
	primePrintCnt := uint32(0)
	number := uint64(from)
	if number <= 2 {
		out <- uint64(2)
		primePrintCnt++
		number = 3
	} else if number%2 == 0 {
		number++
	}
MainLoop:
	for ; number <= uint64(to); number += 2 {
		select {
		case _, isOpen := <-in:
			if !isOpen {
				status = "ABORTED"
				break MainLoop
			}
		default:
			if isPrime(number) {
				out <- number
				primePrintCnt++
			}
		}
	}
	diag <- fmt.Sprintf("%s: Listed %d primes", status, primePrintCnt)
}

// Returns the smallest / first prime factor for the supplied 'number'.
// If 'MaxPrime' is the largest internal prime generated by this package and
// if 'number' > 'MaxPrime' ^ 2, then the method returns '0' (ZERO) indicating
// that the supplied 'number' could not be tested reliably
func GetFirstPrimeFactor(number uint64) uint64 {
	maxPrime := uint64(arrPrimes[maxPrimeCount-1])
	if number > maxPrime*maxPrime {
		return 0 // cannot reliably test
	}

	stopLimit := uint32(math.Trunc(math.Sqrt(float64(number))))
	for i := 0; i < maxPrimeCount && arrPrimes[i] <= stopLimit; i++ {
		if number%uint64(arrPrimes[i]) == 0 {
			return uint64(arrPrimes[i])
		}
	}
	return number // number is prime
}
