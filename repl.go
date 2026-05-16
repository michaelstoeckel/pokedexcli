package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/michaelstoeckel/pokedexcli/internal/pokeapi"
	"github.com/michaelstoeckel/pokedexcli/internal/pokecache"
)

type config struct {
	next string
	prev string
}

type cliCommand struct {
	name        string
	description string
	callback    func(conf *config, cache *pokecache.Cache) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			description: " displays the names of 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			description: "displays the names of previous 20 location areas",
			callback:    commandMapB,
		},
	}
}

func startRepl() error {
	scanner := bufio.NewScanner(os.Stdin)
	conf := &config{next: "", prev: ""}
	cache := pokecache.NewCache(30 * time.Second)

	// repl loop
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()

		words := cleanInput(input)
		if len(words) <= 0 {
			continue
		}
		command := words[0]

		handler := getCommands()[command].callback
		if handler == nil {
			fmt.Println("unknown command")
		} else {
			err := handler(conf, cache)
			if err != nil {
				fmt.Printf("%v", err)
			}
		}
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func commandExit(conf *config, cache *pokecache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config, cache *pokecache.Cache) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for name, command := range getCommands() {
		fmt.Printf("%s: %s\n", name, command.description)
	}
	return nil
}

func commandMap(conf *config, cache *pokecache.Cache) error {
	if conf.next == "" {
		conf.next = "https://pokeapi.co/api/v2/location-area/"
	}

	locations, err := pokeapi.GetLocations(conf.next, cache)
	if err != nil {
		return err
	}

	updateConfig(conf, locations)
	printLocations(locations)

	return nil
}

func commandMapB(conf *config, cache *pokecache.Cache) error {

	if conf.prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locations, err := pokeapi.GetLocations(conf.prev, cache)
	if err != nil {
		return err
	}

	updateConfig(conf, locations)
	printLocations(locations)

	return nil
}

func updateConfig(conf *config, locations pokeapi.Locations) {
	conf.next = locations.Next
	conf.prev = locations.Previous
}

func printLocations(locations pokeapi.Locations) {
	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}
}
