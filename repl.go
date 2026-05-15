package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/michaelstoeckel/pokedexcli/internal/pokeapi"
)

type config struct {
	next string
	prev string
}

type cliCommand struct {
	name        string
	description string
	callback    func(conf *config) error
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
	conf := config{next: "", prev: ""}

	// repl loop
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()

		words := cleanInput(input)
		command := words[0]

		handler := getCommands()[command].callback
		if handler == nil {
			fmt.Println("unknown command")
		} else {
			err := handler(&conf)
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

func commandExit(conf *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(conf *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for name, command := range getCommands() {
		fmt.Printf("%s: %s\n", name, command.description)
	}
	return nil
}

func commandMap(conf *config) error {
	if conf.next == "" {
		conf.next = "https://pokeapi.co/api/v2/location-area/"
	}

	loca, err := pokeapi.GetLocations(conf.next)
	if err != nil {
		return err
	}

	updateConfig(conf, loca)
	printLocations(loca)

	return nil
}

func commandMapB(conf *config) error {

	if conf.prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	loca, err := pokeapi.GetLocations(conf.prev)
	if err != nil {
		return err
	}

	updateConfig(conf, loca)

	printLocations(loca)

	return nil
}

func updateConfig(conf *config, loca pokeapi.Locations) {
	conf.next = loca.Next
	conf.prev = loca.Previous
}

func printLocations(loca pokeapi.Locations) {
	for _, loc := range loca.Results {
		fmt.Println(loc.Name)
	}
}
