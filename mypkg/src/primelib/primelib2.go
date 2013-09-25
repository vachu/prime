package primelib

import "errors"

const primeInterval = 1000
const maxPrimeLimit = 1000 * 1000 // All primes <= 1 million

// The main Prime Repository.  The members are private and are
// accessible only through the associated methods
type Repo struct {
	primes     map[uint64][]uint64
	primeCount uint32
	maxPrime   uint64
}

var PrimeRepo Repo

func init() {
	PrimeRepo.primes = make(map[uint64][]uint64)
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
		primes, _ := GeneratePrimesBetween(r.maxPrime+2, limit)
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

func (r *Repo) storePrimes(primes []uint64) {
	primesLen := len(primes)
	if primesLen == 0 {
		return
	}

	var sliceStart uint32
	quotient := primes[0] / primeInterval
	for i := 0; i < primesLen; i++ {
		if q := primes[i] / primeInterval; q > quotient {
			if _, ok := r.primes[quotient]; !ok {
				r.primes[quotient] = make([]uint64, 0, primeInterval)
			}
			r.storePrimeSlice(primes[sliceStart:i], r.primes[quotient])
			quotient = q
			sliceStart = uint32(i)
		}
	}

	if _, ok := r.primes[quotient]; !ok {
		r.primes[quotient] = make([]uint64, 0, primeInterval)
	}
	r.storePrimeSlice(primes[sliceStart:primesLen], r.primes[quotient])
}

func (r *Repo) storePrimeSlice(src, dst []uint64) (cnt uint32) {

	return
}
