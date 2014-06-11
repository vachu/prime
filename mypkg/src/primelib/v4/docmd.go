package primelib

import (
	"fmt"
	"strconv"
	"strings"
)

const helpMsg = `
list <n>:
    list out the first 'n' primes; n > 0 && < 2^32 - 1

list_between <m> <n>:
    list out all the primes between the numbers 'm' and 'n'
    where 'm' <= 'n' and ('m', 'n') > 0 && < 2^32 - 1

test <n1> [<n2> ...]:
    test each of <n1>, <n2> ... for primality and print
    human-readable results
`

type myChan chan interface{}
type myFunc func([]string) (myChan, myChan, myChan)

var cmd2FuncMap map[string]myFunc

func init() {
	cmd2FuncMap = make(map[string]myFunc)
	cmd2FuncMap["list"] = doListPrimes
	cmd2FuncMap["list_between"] = doListPrimesBetween
	cmd2FuncMap["test"] = doTest
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

func doTest(args []string) (in, out, diag myChan) {
	diag = make(chan interface{}, 1)
	defer close(diag)

	if len(args) == 0 {
		diag <- "ERROR: nothing to test"
		return
	}

	in, out, diag = makeChanTrio(10 * 1024)
	go func() {
		defer close(out)
		defer close(diag)

		for _, arg := range args {
			var num uint64
			_, errFmt := fmt.Sscanf(arg, "%d", &num)
			if errFmt == nil {
				result := fmt.Sprintf("%d: ", num)
				primeFactor := GetFirstPrimeFactor(num)
				switch {
				case primeFactor == 0:
					result += "cannot test reliably"
				case primeFactor == num:
					result += "PRIME"
				default:
					result += fmt.Sprintf("divisible by %d", primeFactor)
				}
				out <- result
			} else {
				out <- fmt.Sprintf("ERROR: %s - not a uint32 number", arg)
			}
		}
	}()
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
