#!/usr/bin/perl

use strict;
use English;

my ($prime, $lastPrime) = (undef, 2);
while (defined($prime = <>)) {
	chomp($prime);
	print "$lastPrime, $prime\n" if ($prime - $lastPrime == 2);
	$lastPrime = $prime;
}
