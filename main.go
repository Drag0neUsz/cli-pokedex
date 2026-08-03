package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		_ = scanner.Scan()
		if scanner.Err() != nil {
			break
		}
		input := scanner.Text()
		command, ok := getCommands()[input]
		if !ok {
			fmt.Print("Unknown command\n")
			continue
		}
		command.callback()
	}
}
