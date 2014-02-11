# prime
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

## Motivation
I want to learn Golang and I have a student's fascination for Prime Numbers.
I also want to learn how to create a simple ReST-ful WebService and the new
WebSockets technology.  Hence I devised the project solely for my own learning
purposes.  In the process, I also want get a good grip on git.

The "distant" future plan is to come up with a HTML5 client that consumes the
WebService as well as the WebSocket server

## Enhancement Plans
The PrimeLib would be served by:
<ol>
	<li>a simple ReST WebService</li>
	<li>a WebSocket server</li>
</ol>

Since there is no concept of shared library / DLL in go (as of 04-Feb-2014),
the immediate future plan is to wrap such librarires into a WebSocket server
so that other apps could "reuse" them at runtime.  In that light, the currently
planned ReST WebService would also "reuse" the WebSocket-served library

## More (Gory) Details


:-x-x-x-:
