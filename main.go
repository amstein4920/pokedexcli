package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		err := getCommands()[input].callback()
		if err != nil {
			fmt.Printf("error with callback: %v", err)
		}
	}
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display previous 20 locations",
			callback:    commandMapB,
		},
	}
}

func commandHelp() error {
	var commandsSorted []cliCommand

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, value := range getCommands() {
		commandsSorted = append(commandsSorted, value)
	}
	sort.Slice(commandsSorted, func(i, j int) bool {
		return i < j
	})

	for _, command := range commandsSorted {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func commandMap() error {
	fmt.Println("Next")
	return nil
}

func commandMapB() error {
	fmt.Println("Previous")
	return nil
}
