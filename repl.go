package main

import (
	"bufio"
	"fmt"
	"math/rand"
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
	Response pokeapi.PokemapResponse
	Pokedex  map[string]pokeapi.Pokemon
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
			description: "explores a map",
			callback:    p.commandExplore,
		},
		"catch": {
			description: "catches a pokemon",
			callback:    p.commandCatch,
		},
		"inspect": {
			description: "inspects a pokemon",
			callback:    p.commandInspect,
		},
		"pokedex": {
			description: "shows the content of your pokedex",
			callback:    p.commandPokedex,
		},
	}
}

func startRepl() error {
	scanner := bufio.NewScanner(os.Stdin)
	conf := config{next: "", prev: ""}
	pokeapi := Pokeapi{
		Cache:    *pokecache.NewCache(30 * time.Second),
		Response: pokeapi.PokemapResponse{},
		Pokedex:  map[string]pokeapi.Pokemon{},
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
	p.Response, err = pokeapi.GetPokemapResponse(url, &p.Cache)
	if err != nil {
		return err
	}

	for _, encounter := range p.Response.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}

func (p *Pokeapi) commandCatch(args []string) error {
	url := "https://pokeapi.co/api/v2/pokemon/"
	if len(args) == 0 {
		return fmt.Errorf("please specify a pokemon\n")
	}
	name := args[0]
	url += name
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	pokemon, err := pokeapi.GetPokemon(url, &p.Cache)
	if err != nil {
		return err
	}

	catchVal := rand.Intn(500)
	fmt.Printf("%s %v/%v\n", pokemon.Name, pokemon.BaseExperience, catchVal)

	if catchVal < pokemon.BaseExperience {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	p.Pokedex[pokemon.Name] = pokemon
	fmt.Printf("%s was caught!\n", pokemon.Name)

	return nil
}

func (p *Pokeapi) commandInspect(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please specify a pokemon\n")
	}
	name := args[0]

	pokemon, ok := p.Pokedex[name]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("Name:   %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	var types pokeapi.Types
	for _, types = range pokemon.Types {
		fmt.Printf("  -%v\n", types.Type.Name)
	}

	return nil
}

func (p *Pokeapi) commandPokedex(args []string) error {
	fmt.Println("Your pokedex:")
	for _, pokemon := range p.Pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
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
