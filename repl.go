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
	callback    func(args []string) error
}

type Pokeapi struct {
	Cache    pokecache.Cache
	Response pokeapi.PokemonResponse
	config   config
}

func getCommands(p *Pokeapi) map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			description: "Exit the Pokedex",
			callback:    p.commandExit,
		},
		"help": {
			description: "Displays a help message",
			callback:    p.commandHelp,
		},
		"map": {
			description: " displays the names of 20 location areas",
			callback:    p.commandMap,
		},
		"mapb": {
			description: "displays the names of previous 20 location areas",
			callback:    p.commandMapB,
		},
		"explore": {
			description: "what?",
			callback:    p.commandExplore,
		},
	}
}

func startRepl() error {
	scanner := bufio.NewScanner(os.Stdin)
	conf := config{next: "", prev: ""}
	pokeapi := Pokeapi{
		Cache:    *pokecache.NewCache(30 * time.Second),
		Response: pokeapi.PokemonResponse{},
		config:   conf,
	}

	// repl loop
	for {
		fmt.Print("Pokedex > ")

		scanner.Scan()
		input := scanner.Text()

		words := cleanInput(input)
		if len(words) <= 0 {
			continue
		}
		commandName := words[0]
		args := words[1:]
		commands := getCommands(&pokeapi)

		if command, exists := commands[commandName]; exists {
			err := command.callback(args)
			if err != nil {
				fmt.Printf("Error: %v", err)
			}
		}
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func (p *Pokeapi) commandExit(args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func (p *Pokeapi) commandHelp(args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for name, command := range getCommands(p) {
		fmt.Printf("%s: %s\n", name, command.description)
	}
	return nil
}

func (p *Pokeapi) commandMap(args []string) error {
	if p.config.next == "" {
		p.config.next = "https://pokeapi.co/api/v2/location-area/"
	}

	locations, err := pokeapi.GetLocations(p.config.next, &p.Cache)
	if err != nil {
		return err
	}

	updateConfig(&p.config, locations)
	printLocations(locations)

	return nil
}

func (p *Pokeapi) commandMapB(args []string) error {

	if p.config.prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locations, err := pokeapi.GetLocations(p.config.prev, &p.Cache)
	if err != nil {
		return err
	}

	updateConfig(&p.config, locations)
	printLocations(locations)

	return nil
}

func (p *Pokeapi) commandExplore(args []string) error {
	url := "https://pokeapi.co/api/v2/location-area/"
	if len(args) == 0 {
		return fmt.Errorf("please specify an area\n")
	}
	url += args[0]
	fmt.Printf("Exploring %s ...\n", args[0])
	var err error
	p.Response, err = pokeapi.GetPokemonResponse(url, &p.Cache)
	if err != nil {
		return err
	}

	for _, encounter := range p.Response.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

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
