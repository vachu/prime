package primelib

import "math"

const MaxPrimeCount = 1000 * 1000

var arrPrimes [MaxPrimeCount]uint32

func init() {
	arrPrimes[0] = 2
	primeCount := 1
	for number := uint32(3); primeCount < MaxPrimeCount; number += 2 {
		if isPrime(number) {
			arrPrimes[primeCount] = number
			primeCount++
		}
	}
}

func isPrime(number uint32) bool {
	isComposite := false
	stopLimit := uint32(math.Trunc(math.Sqrt(float64(number))))
	for i := 0; !isComposite && arrPrimes[i] <= stopLimit; i++ {
		isComposite = (number%arrPrimes[i] == 0)
	}
	return !isComposite
}

func ListPrimes(out chan uint32, cnt uint32) {
	for i := uint32(0); i < cnt && i < MaxPrimeCount; i++ {
		out <- arrPrimes[i]
	}

	number := arrPrimes[MaxPrimeCount-1] + 2
	for i := uint32(MaxPrimeCount); i < cnt && number <= math.MaxUint32; number += 2 {
		if isPrime(number) {
			out <- number
			i++
		}
	}

	close(out)
}

func ListPrimesBetween(out chan uint32, from, to uint32) {
	if to > 0 && from <= to {
		number := from
		if number <= arrPrimes[MaxPrimeCount-1] {
			// print from the internal prime array
			for i := 0; i < MaxPrimeCount; i++ {
				if arrPrimes[i] >= from && arrPrimes[i] <= to {
					out <- arrPrimes[i]
				}
			}
			number = arrPrimes[MaxPrimeCount-1] + 2
		} else if number%2 == 0 {
			number++
		}
		for ; number >= from && number <= to; number += 2 {
			if isPrime(number) {
				out <- number
			}
		}
	}
	close(out)
}

func GetFirstPrimeFactor(number uint64) uint64 {
	if number > math.MaxUint32 {
		return 0 // cannot use a number > 2^32-1
	}

	stopLimit := uint32(math.Trunc(math.Sqrt(float64(number))))
	for i := 0; i < MaxPrimeCount && arrPrimes[i] <= stopLimit; i++ {
		if number%uint64(arrPrimes[i]) == 0 {
			return uint64(arrPrimes[i])
		}
	}
	return number // number is prime
}
