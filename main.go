package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	for {
		fmt.Print("Pokedex > ")
		_ = scanner.Scan()
		if scanner.Err() != nil {
			break
		}
		input := scanner.Text()
		parts := cleanInput(input)
		command, ok := commands[parts[0]]
		if !ok {
			fmt.Print("Unknown command\n")
			continue
		}
		err := command.callback(command.callbackConfig, parts[1:]...)
		if err != nil {
			fmt.Printf("Error executing command: %v\n", err)
		}
	}
}
