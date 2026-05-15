package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type config struct {
	next string
	prev string
}

type location struct {
	// Capitalize Name and Url
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locations struct {
	// Capitalize all fields
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []location `json:"results"`
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

	res, err := http.Get(conf.next)
	if err != nil {
		return err // Return the error instead of killing the app
	}
	defer res.Body.Close() // Ensure closure happens

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	loca := locations{}
	err = json.Unmarshal(body, &loca)
	if err != nil {
		return err
	}

	// Now you can update your config for the "next" page
	conf.next = loca.Next
	conf.prev = loca.Previous

	for _, loc := range loca.Results {
		fmt.Println(loc.Name)
	}

	return nil
}

func commandMapB(conf *config) error {

	if conf.prev == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	res, err := http.Get(conf.prev)
	if err != nil {
		return err // Return the error instead of killing the app
	}
	defer res.Body.Close() // Ensure closure happens

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	loca := locations{}
	err = json.Unmarshal(body, &loca)
	if err != nil {
		return err
	}

	// Now you can update your config for the "next" page
	conf.next = loca.Next
	conf.prev = loca.Previous

	for _, loc := range loca.Results {
		fmt.Println(loc.Name)
	}

	return nil
}
