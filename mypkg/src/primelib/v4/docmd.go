package primelib

import (
	"fmt"
	"strconv"
	"strings"
)

const helpMsg = `Available Commands (case-sensitive):

ListPrimes n: list out the first 'n' primes; n > 0 && < 2^32 - 1

ListPrimesBetween m n: list out all the primes between the numbers 'm' and 'n'
          where 'm' <= 'n' and ('m', 'n') > 0 && < 2^32 - 1
`

type myFunc func([]string) (myChan, myChan, myChan)

var cmd2FuncMap map[string]myFunc

func init() {
	cmd2FuncMap = make(map[string]myFunc)
	cmd2FuncMap["ListPrimes"] = doListPrimes
	cmd2FuncMap["ListPrimesBetween"] = doListPrimesBetween
	cmd2FuncMap["?"] = doPrintHelp
}

func DoCmd(cmdLine string) (in, out, diag myChan) {
	diag = make(chan interface{}, 1)
	defer close(diag)

	if len(cmdLine) > 0 {
		fields := strings.Fields(cmdLine)
		if fn, ok := cmd2FuncMap[fields[0]]; ok {
			in, out, diag = fn(fields[1:])
		} else {
			diag <- fmt.Sprintf("ERROR: illegal command - (%s)", cmdLine)
		}
	} else {
		diag <- fmt.Sprintf("ERROR: nil / empty cmdLine (%s)", cmdLine)
	}
	return
}

func doPrintHelp(args []string) (in, out, diag myChan) {
	diag = make(chan interface{}, 1)
	diag <- helpMsg
	close(diag)
	return
}

func doListPrimes(args []string) (in, out, diag myChan) {
	diag = make(chan interface{}, 1)
	defer close(diag)

	if len(args) != 1 {
		diag <- "ERROR: illegal arg count - expected 1 arg"
		return
	}

	var cnt uint32
	var e error
	if cnt, e = getUint32Val(args[0]); e != nil {
		diag <- e.Error()
		return
	}

	in, out, diag = ListPrimes(cnt)
	return
}

func doListPrimesBetween(args []string) (in, out, diag myChan) {
	diag = make(chan interface{}, 1)
	defer close(diag)

	if len(args) != 2 {
		diag <- "ERROR: illegal arg count - expected 2 args"
		return
	}

	var from, to uint32
	var e error
	if from, e = getUint32Val(args[0]); e != nil {
		diag <- fmt.Sprintf("%s (arg1)", e.Error())
		return
	}
	if to, e = getUint32Val(args[1]); e != nil {
		diag <- fmt.Sprintf("%s (arg2)", e.Error())
		return
	}

	in, out, diag = ListPrimesBetween(from, to)
	return
}

func getUint32Val(arg string) (number uint32, e error) {
	str := strings.Trim(arg, " \r\n\t\v\f")
	if len(str) == 0 {
		e = fmt.Errorf("ERROR: empty / nil arg value (%s)", arg)
	} else {
		var n64 uint64
		n64, e = strconv.ParseUint(str, 0, 32)
		if e == nil {
			number = uint32(n64)
		} else {
			e = fmt.Errorf("ERROR: illegal uint32 value - %s", e.Error())
		}
	}
	return
}
