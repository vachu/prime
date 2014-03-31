# prime
My first, major program in Google's Go Programming Language!

Right now (as of 31-Aug-2013 IST), this program provides for a simple
console-based UI with which we could:<br/>
	<ol>
	<li>generate 'n' primes</li>
	<li>load primes from text file - 1 prime number per line<br/>
		<em>NOTE: this feature is only in the old 'prime.go' and it is assumed that this text file consists of sequential
		primes.  No checks / validation are made.  Hence it is possible
		to feed WRONG primes as well.</em></li>
	<li>test a given number via STDIN for primality</li>
	</ol>

The programs 'prime.go' and 'prime2.go' are implemented as a regular Unix filter.
'prime_wss.go' is a simple Websocket Server app that serves the 'primelib/v3'
package.  It takes simple text-based commands like 'LIST' and 'TEST' to provide
for the following functionlities:
	<ol>
		<li>List the available commands to the Websocket client</li>
		<li>List the first 'n' primes where 'n' >=0 && < MAX_INT32</li>
		<li>List all the primes between 'm' and 'n' where 'm' <= 'n' &&
		('m', 'n') >= 0 && < MAX_INT32</li>
		<li>Test the given number(s) from primality; the result of
		each check is output to the Websocket Client</li>
	</ol>

The commands could be given to the running Websocket Server App through the
interactive (websocket) Client App 'prime_wsc.go'

All the programs in this project are console-based apps.

## Motivation
I want to learn Golang and I have a student's fascination for Prime Numbers.
I also want to learn how to create a simple ReST-ful WebService and the new
WebSockets technology.  Hence I devised the project solely for my own learning
purposes.  In the process, I also want get a good grip on git.

The "distant" future plan is to come up with a HTML5 client that consumes the
WebService as well as the WebSocket server

## More (Gory) Details
The source has been restructured.  The old and now almost defunct 'prime.go' uses the 'primelib/v1' package and the latest 'prime2.go' uses the 'primelib/v3'.  'primelib/defunt' was supposed to be 'primelib/v2' but something was not alright with it

### Primelib/v3
This package is greatly simplified from v1.  It uses an internal array of a million primes (generated @ startup/init) to print primes and test input numbers for primality

This package is now served through a simple Websocket server program 'prime_wss.go' that uses the 3rd-party package 'code.google.com/p/go.net/websocket'.  'prime_wsc.go' is the Websocket client program that connects to the above Websocket Server program.

## Enhancement Plans
The PrimeLib would be served by:
<ol>
	<li>a simple ReST WebService</li>
	<li>a WebSocket server (Update: a basic version is now in place)</li>
</ol>

Since there is no concept of shared library / DLL in go (as of 04-Feb-2014),
the 'primelib/v3' package is wrapped into a WebSocket server so that other apps
could "reuse" them at runtime.  In that light, the currently planned 
ReST WebService would also "reuse" the WebSocket-served library

:-x-x-x-:
