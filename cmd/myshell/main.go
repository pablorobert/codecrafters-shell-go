package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		// Wait for user input
		str, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			panic(err)
		}
		str = strings.Replace(str, "\n", "", -1)
		ParseCommand(str)
	}
}

func ParseCommand(str string) {
	parts := strings.Split(str, " ")
	switch parts[0] {
	case "exit":
		var exitCode int
		exitCode, err := strconv.Atoi(parts[1])
		if err != nil {
			exitCode = 0
		}
		os.Exit(exitCode)
	case "echo":
		output := strings.Join(parts[1:], " ")
		fmt.Printf("%s\n", output)
	default:
		fmt.Printf("%s: command not found\n", parts[0])
	}
}
