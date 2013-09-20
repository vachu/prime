package primelib

import "errors"

type Repo struct {
	primes     map[uint32][]uint64
	primeCount uint32
	maxPrime   uint64
}

var primeRepo Repo

const PRIME_INTERVAL = 1000
const MAX_PRIME_LIMIT = 1000 * 1000 * 1000 // All primes <= 1 billion

func init() {
	primeRepo.primes = make(map[uint32][]uint64)
}

// Generates first 'n' primes where 'n' <= 'limit' and stores them
// in the internal prime repository.
//
// If the internal prime repository already contains the required list
// of primes and there is no need to generate additional primes, then a
// count of ZERO is returned.
//
// If some additional primes are needed to be generated upto 'limit',
// then the count of only those newly generated primes is returned.
func (r *Repo) GeneratePrimesUpto(limit uint64) (cnt uint32) {
	switch {
	case limit <= 1 || r.maxPrime >= limit:
		cnt = 0
	default:
		primes, _ := GeneratePrimesBetween(r.maxPrime+1, limit)
		cnt = uint32(len(primes))
	}

	return
}

// Generates primes between the numbers 'lo' and 'hi' where 'lo' < 'hi'
func GeneratePrimesBetween(lo, hi uint64) (primes []uint64, err error) {
	if lo >= hi {
		err = errors.New("lo >= hi")
		return
	}

	return
}
