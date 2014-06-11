package main

import (
	"bufio"
	"fmt"
	"os"
	"primelib/v4"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			processCmdline(arg)
		}
	} else {
		for cmdLine := getLine(); len(cmdLine) > 0; cmdLine = getLine() {
			processCmdline(cmdLine)
		}
	}

}

func getLine() string {
	fmt.Print("Enter command: ")
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.Trim(line, " \t\v\f\r\n")
	line = strings.ToLower(line)
	return line
}

func processCmdline(cmdLine string) {
	_, out, diag := primelib.DoCmd(cmdLine)
	if out != nil {
		for output := range out {
			fmt.Println(output)
		}
	}
	for diagOutput := range diag {
		fmt.Fprintln(os.Stderr, diagOutput)
	}
}
