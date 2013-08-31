prime
=====

My first, major program in Google's Go Programming Language!

Right now (as of 31-Aug-2013 IST), this program provides for a simple
console-based UI with which we could:<br/>
	<ol>
	<li>generate 'n' primes</li>
	<li>load primes from text file - 1 prime number per line<br/>
		<em>NOTE: it is assumed that this text file consists of sequential
		primes.  No checks / validation are made.  Hence it is possible
		to feed WRONG primes as well.</em></li>
	<li>test a given number via STDIN for primality</li>
	</ol>
The whole program is implemented as a regular Unix filter.

Future Plans
~~~~~~~~~~~~
The Primelib package / library will be converted to a suitable Service.
I am planning to use WebSockets instead of the current WebService 
mechanisms.

-x-x-x-
