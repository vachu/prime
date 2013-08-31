prime
=====

My first, major program in Google's Go Programming Language!

Right now (as of 31-Aug-2013 IST), this program provides for a simple
console-based UI with which we could:
	1. generate 'n' primes
	2. load primes from text file - 1 prime number per line
		NOTE: it is assumed that this text file consists of sequential
		primes.  No checks / validation are made.  Hence it is possible
		to feed WRONG primes as well.
	3. test a given number <via STDIN> for primality

The whole program is implemented as a regular Unix filter.

Future Plans
~~~~~~~~~~~~
The Primelib package / library will be converted to a suitable Service.
I am planning to use WebSockets instead of the current WebService 
mechanisms.

-x-x-x-
