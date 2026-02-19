package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"slices"
	"strconv"
	"strings"
)

var bultins []string = []string{"type", "echo", "exit", "cat"}

func main() {
	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, os.Interrupt)
    go func() {		
		<-sigTerm

		fmt.Println("Received SIGTERM signal")
		os.Exit(1)
    }()
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

func handleCat(parts []string) {
	result := ""
	if len(parts) < 2 {
		fmt.Println("cat: missing file operand")
		return
	}
	for i := 1; i < len(parts); i++ {
		fileName := parseSingleQuotes(parts[i])
		data, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println("Error while reading file")
			return
		}
		result = result + fmt.Sprintf("%s\n", string(data))
	}
	result += "\n"
	fmt.Print(result)
}

func parseSingleQuotes(str string) string {
			var quoteCount int
			var inSpace bool
			var returnStr strings.Builder
 			for _, runeValue := range str {
				if runeValue == '\'' {
					quoteCount++
					continue
				}
				if (runeValue == '\'') {
					quoteCount--
					continue
				}
				if quoteCount % 2 == 0 && runeValue == ' ' && inSpace {
					continue
				}
				if runeValue == ' ' {
					inSpace = true
				} else {
					inSpace = false
				}
				
				returnStr.WriteString(string(runeValue))
			}

			if quoteCount % 2 != 0 {
				panic("Unmatched single quotes")
			}

			return returnStr.String()
		}
		


func ParseCommand(str string) {
	parts := strings.Split(str, " ")
	switch parts[0] {
	case "exit":
		if len(parts) < 2 {
			os.Exit(0)
		}
		var exitCode int
		exitCode, err := strconv.Atoi(parts[1])
		if err != nil {
			exitCode = 0
		}
		os.Exit(exitCode)
	case "echo":
		output := strings.Join(parts[1:], " ") 
		output = parseSingleQuotes(output)
		fmt.Printf("%s\n", output)
	case "cd":
		if parts[1] == "~" {
			parts[1] = os.Getenv("HOME")
		}
		_, err := os.Stat(parts[1])
		if os.IsNotExist(err) {
			fmt.Printf("%s: No such file or directory\n", parts[1])
			return
		}
		os.Chdir(parts[1])
	case "pwd":
		output, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
	case "type":
		cmd := strings.TrimSpace(parts[1])
		idx := slices.IndexFunc(bultins, func(c string) bool { return c == cmd })

		if idx == -1 {
			GetPath(cmd)
		} else {
			fmt.Fprintf(os.Stdout, "%v is a shell builtin\n", cmd)
		}
	case "cat":
		/*if len(parts) < 2 {
			fmt.Println("cat: missing file operand")
			return
		}
		cmd := strings.TrimSpace(parts[1])
		idx := slices.IndexFunc(bultins, func(c string) bool { return c == cmd })

		if idx == -1 {
			GetPath(cmd)
		} else {
			fmt.Fprintf(os.Stdout, "%v is a shell builtin\n", cmd)
		}*/
		handleCat(parts)
	default:
		fmt.Println("cat")
		err := Execute(str)
		if err != nil {
			fmt.Printf("%s: command not found\n", parts[0])
		}
	}
}

func Execute(cmd string) error {
	program := strings.Split(cmd, " ")
	exe := program[0]

	_, err := os.Stat(exe)
	if err != nil {
		return err
	}

	args := program[1:]
	output, err := exec.Command(exe, args...).Output()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "%v", string(output))
	return nil
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
