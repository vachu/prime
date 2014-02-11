package primelib

import (
	"fmt"
	"io"
	"math"
)

const MaxPrimeCount = 6542 // total no. of unsigned 16-bit primes
var arrPrimes [MaxPrimeCount]uint16

func init() {
	arrPrimes[0] = 2
	primeCount := 1

	for number := uint32(3); number <= math.MaxUint16; number += 2 {
		if isPrime(number) {
			arrPrimes[primeCount] = uint16(number)
			primeCount++
		}
	}
}

func isPrime(number uint32) bool {
	isComposite := false
	stopLimit := uint16(math.Trunc(math.Sqrt(float64(number))))
	for i := 0; !isComposite && arrPrimes[i] <= stopLimit; i++ {
		isComposite = (number%uint32(arrPrimes[i]) == 0)
	}
	return !isComposite
}

func WritePrimes(w io.Writer, cnt uint32, sep string) uint32 {
	for i := uint32(0); i < cnt && i < MaxPrimeCount; i++ {
		fmt.Fprintf(w, "%d%s", arrPrimes[i], sep)
	}

	number := uint32(arrPrimes[MaxPrimeCount-1]) + 2
	for i := uint32(MaxPrimeCount); i < cnt; number += 2 {
		if isPrime(number) {
			fmt.Fprintf(w, "%d%s", number, sep)
			i++
		}
	}

	return cnt
}

func WritePrimesBetween(w io.Writer, from, to uint32, sep string) (cnt uint32) {
	cnt = 0
	if to == 0 || from > to {
		return
	}

	number := from
	if number%2 == 0 {
		number++
	}

	for ; number >= from && number <= to; number += 2 {
		if isPrime(number) {
			fmt.Fprintf(w, "%d%s", number, sep)
		}
	}

	return
}

func GetFirstPrimeFactor(number uint64) uint64 {
	if number > math.MaxUint32 {
		return 0 // cannot use a number > 2^32-1
	}

	stopLimit := uint16(math.Trunc(math.Sqrt(float64(number))))
	for i := 0; i < MaxPrimeCount && arrPrimes[i] <= stopLimit; i++ {
		if number%uint64(arrPrimes[i]) == 0 {
			return uint64(arrPrimes[i])
		}
	}
	return number // number is prime
}
