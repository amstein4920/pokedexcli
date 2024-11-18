package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/amstein4920/pokedexcli/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

type config struct {
	pokeClient          pokeapi.Client
	caughtPokemon       map[string]pokeapi.Pokemon
	nextLocationUrl     *string
	previousLocationUrl *string
}

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		caughtPokemon: map[string]pokeapi.Pokemon{},
		pokeClient:    pokeClient,
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		words := strings.Fields(input)

		arguments := []string{}
		if len(words) > 1 {
			arguments = words[1:]
		}

		command, exists := getCommands()[words[0]]
		if exists {
			err := command.callback(cfg, arguments...)
			if err != nil {
				fmt.Println(err)
			}
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
		"explore": {
			name:        "explore",
			description: "Explore given location and display Pokemon within it",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch given Pokemon and add to your personal Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Show information for the given Pokemon if the Pokemon has been previously caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show names of all caught Pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandHelp(config *config, args ...string) error {
	var commandsSorted []cliCommand

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for _, value := range getCommands() {
		commandsSorted = append(commandsSorted, value)
	}
	sort.Slice(commandsSorted, func(i, j int) bool {
		return commandsSorted[i].name < commandsSorted[j].name
	})

	for _, command := range commandsSorted {
		fmt.Printf("%v: %v\n", command.name, command.description)
	}
	return nil
}

func commandExit(config *config, args ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *config, args ...string) error {
	locations, err := config.pokeClient.ListLocations(config.nextLocationUrl)
	if err != nil {
		return err
	}
	config.nextLocationUrl = locations.Next
	config.previousLocationUrl = locations.Previous

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapB(config *config, args ...string) error {
	locations, err := config.pokeClient.ListLocations(config.previousLocationUrl)
	if err != nil {
		return err
	}
	config.nextLocationUrl = locations.Next
	config.previousLocationUrl = locations.Previous

	if config.previousLocationUrl == nil {
		return errors.New("first page, can not go back further")
	}

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(config *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("no explore argument provided")
	}

	location, err := config.pokeClient.ExploreLocation(args[0])
	if err != nil {
		return err
	}

	for _, encounter := range location.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(config *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("incorrect number of arguments provided")
	}

	pokemon, err := config.pokeClient.GetPokemon(args[0])
	if err != nil {
		return err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	is_caught := r.Intn(300) >= pokemon.BaseExperience

	if is_caught {
		fmt.Println(pokemon.Name + " has been caught and added to your Pokedex!")
		config.caughtPokemon[pokemon.Name] = pokemon
	} else {
		fmt.Println(pokemon.Name + " was not caught!")
	}

	return nil
}

func commandInspect(config *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("incorrect number of arguments provided")
	}

	pokemon, ok := config.caughtPokemon[args[0]]
	if !ok {
		fmt.Println("Pokemon provided has not been caught")
		return nil
	}

	printInfo(pokemon)
	return nil
}

func printInfo(pokemon pokeapi.Pokemon) {
	fmt.Println("Name: " + pokemon.Name)
	fmt.Println("Height: " + fmt.Sprint(pokemon.Height))
	fmt.Println("Weight: " + fmt.Sprint(pokemon.Weight))

	fmt.Println("Stats:")
	for _, val := range pokemon.Stats {
		fmt.Println("-" + val.Stat.Name + ": " + fmt.Sprint(val.BaseStat))
	}

	fmt.Println("Types:")
	for _, val := range pokemon.Types {
		fmt.Println("-" + val.Type.Name)
	}
}

func commandPokedex(config *config, args ...string) error {
	if len(config.caughtPokemon) < 1 {
		fmt.Println("Pokedex is Empty")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.caughtPokemon {
		fmt.Println("-" + pokemon.Name)
	}
	return nil
}
