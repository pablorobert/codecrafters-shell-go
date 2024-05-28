package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

var bultins []string = []string{"type", "echo", "exit"}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		// Wait for user input
		str, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			panic(err)
		}
		str = strings.TrimSpace(str)
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
	case "type":
		cmd := strings.TrimSpace(parts[1])
		idx := slices.IndexFunc(bultins, func(c string) bool { return c == cmd })

		if idx == -1 {
			GetPath(cmd)
		} else {
			fmt.Fprintf(os.Stdout, "%v is a shell builtin\n", cmd)
		}

	default:
		fmt.Printf("%s: command not found\n", parts[0])
	}
}

func GetPath(cmd string) {
	path := os.Getenv("PATH")
	parts := strings.Split(path, ":")
	found := false
	var exe string
	for _, p := range parts {
		exe = p + "/" + cmd
		_, err := os.Stat(exe)
		if err != nil {
			continue
		}
		found = true
		break
	}

	if found {
		fmt.Fprintf(os.Stdout, "%v is %v\n", cmd, exe)
	} else {
		fmt.Fprintf(os.Stdout, "%v not found\n", cmd)
	}
}
