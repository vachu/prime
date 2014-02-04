package primelib

import ("fmt"; "math"; "io")

const maxPrimeCount = 6542 // no. of primes <= 2 ^ 16 - 1
var arrPrimes [maxPrimeCount]uint16

func init() {
	arrPrimes[0] = 2
	primeCount := 1

	for number := 3; number <= math.MaxUint16; number += 2 {
		isComposite := false;
		stopLimit := uint16(math.Trunc(math.Sqrt(float64(number))))
		for i := 0; !isComposite && arrPrimes[i] <= stopLimit; i++ {
			isComposite = (number % int(arrPrimes[i]) == 0)
		}

		if (!isComposite) {
			arrPrimes[primeCount] = uint16(number)
			primeCount++
		}
	}
}

func WritePrimes(w io.Writer, cnt uint, sep string) {
	if cnt > maxPrimeCount {
		cnt = maxPrimeCount
	}

	for i := uint(0); i < cnt; i++ {
		fmt.Fprintf(w, "%d%s", arrPrimes[i], sep)
	}
}

func WritePrimesBetween(w io.Writer, from, to uint, sep string) uint {
	return 0
}
