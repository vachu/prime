package primelib

import "math"
import "os"
import "fmt"

const MAX_ARR_ELEM = 1024

var primes map[uint][]uint64
var primeCount uint
var maxPrime uint64

func init() {
	primes = make(map[uint][]uint64)
}

func LoadPrimes(count uint, inputPrimesFile string) (uint, error) {
	primes_file, err := os.Open(inputPrimesFile)
	if err != nil {
		return 0, err
	}
	defer primes_file.Close()

	var prime uint64
	_, err = fmt.Fscanf(primes_file, "%d", &prime)
	for err == nil && primeCount < count {
		quotient := primeCount / MAX_ARR_ELEM
		remainder := primeCount % MAX_ARR_ELEM
		if remainder == 0 {
			primes[quotient] = make([]uint64, MAX_ARR_ELEM)
		}

		primes[quotient][remainder] = prime
		maxPrime = prime
		primeCount++
		_, err = fmt.Fscanf(primes_file, "%d", &prime)
	}

	return primeCount, err
}

func GeneratePrimes(count uint) uint {
	for primeCount < count {
		quotient := primeCount / MAX_ARR_ELEM
		remainder := primeCount % MAX_ARR_ELEM
		if remainder == 0 {
			primes[quotient] = make([]uint64, MAX_ARR_ELEM)
		}

		var candidate uint64
		switch maxPrime {
		case 0:
			candidate = 2
		case 2:
			candidate = 3
		default:
			candidate = maxPrime + 2
			for ; !IsPrime(candidate); candidate += 2 {
			}
		}

		primes[quotient][remainder] = candidate
		maxPrime = candidate
		primeCount++
	}
	return primeCount
}

func GetPrime(index uint) uint64 {
	if index >= primeCount {
		return 0
	}

	quotient := index / MAX_ARR_ELEM
	remainder := index % MAX_ARR_ELEM
	return primes[quotient][remainder]
}

func GetFirstPrimeFactor(number uint64) uint64 {
	if primeCount == 0 {
		panic("No primes generated")
	}

	stopLimit := uint64(math.Trunc(math.Sqrt(float64(number))))
	var i uint
	prime := GetPrime(i)
	for ; prime != 0 && prime <= stopLimit; prime = GetPrime(i) {
		if number%prime == 0 {
			return prime
		}
		i++
	}
	return number
}

func IsPrime(candidate uint64) bool {
	return GetFirstPrimeFactor(candidate) == candidate
}

func GetPrimeCount() uint {
	return primeCount
}

func GetMaxPrime() uint64 {
	return maxPrime
}
