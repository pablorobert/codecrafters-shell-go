package main

import (
	"bufio"
	"fmt"
	"os"
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
		fmt.Printf("%s: command not found\n", str)
	}
}
