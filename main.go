package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	// "./pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type PokeLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// type config struct {
// pokeClient pokeapi.Client

// }

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
	res, err := http.Get("https://pokeapi.co/api/v2/location/")
	if err != nil {
		fmt.Println("Error loading locations")
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading locations")
		return err
	}
	defer res.Body.Close()

	locations := PokeLocations{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println("Error with locations JSON")
		return err
	}
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapB() error {
	fmt.Println("Previous")
	return nil
}
