package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fbb-mk1/pokedex/internal/pkcache"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string) error
}

type Config struct {
	cache                pkcache.Cache
	nextLocationsURL     string
	previousLocationsURL *string
	seenPokemon          map[string]string
}

func getCommands() map[string]CliCommand {
	return map[string]CliCommand{
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
			description: "Show the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previousLocationsURL 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Check for pokemons in a specific area, using area ID or Name",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon that you have seen while exploring",
			callback:    commandExplore,
		},
	}
}

func commandMap(config *Config, p ...string) error {
	local, ok := config.cache.Get(config.nextLocationsURL)
	if !ok {
		l, err := getLocation(config.nextLocationsURL)
		if err != nil {
			return err
		}
		local = l
		config.cache.Add(config.nextLocationsURL, local)
	}
	l, err := getLocationValues(local)
	if err != nil {
		return err
	}
	config.nextLocationsURL = l.Next
	config.previousLocationsURL = l.Previous
	fmt.Println("ID: Area Name")
	for _, loc := range l.Results {
		id := strings.Split(loc.URL, "/")
		line := fmt.Sprintf("%v: %v", id[len(id)-2], loc.Name)
		fmt.Println(line)
	}
	return nil
}

func commandMapb(config *Config, p ...string) error {
	if config.previousLocationsURL == nil {
		return errors.New("at start of list")
	}
	local, ok := config.cache.Get(*config.previousLocationsURL)
	if !ok {
		l, err := getLocation(*config.previousLocationsURL)
		if err != nil {
			return err
		}
		local = l
		config.cache.Add(*config.previousLocationsURL, local)
	}
	l, err := getLocationValues(local)
	if err != nil {
		return err
	}
	config.nextLocationsURL = l.Next
	config.previousLocationsURL = l.Previous
	for _, loc := range l.Results {
		id := strings.Split(loc.URL, "/")
		line := fmt.Sprintf("%v: %v", id[len(id)-2], loc.Name)
		fmt.Println(line)
	}
	return nil
}

func commandHelp(*Config, ...string) error {
	commands := getCommands()
	fmt.Println("========POKEDEX========")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("How to use:")
	fmt.Println("")
	for key := range commands {
		fmt.Printf("%v: %v\n", commands[key].name, commands[key].description)
	}
	return nil
}

func commandExit(*Config, ...string) error {
	fmt.Println("Bye Bye")
	os.Exit(0)
	return nil
}

func commandExplore(config *Config, p ...string) error {
	if p[0] == " " {
		return errors.New("include an area name or ID")
	}
	search := "https://pokeapi.co/api/v2/location-area/" + p[0]
	result, ok := config.cache.Get(search)
	if !ok {
		r, err := getLocation(search)
		if err != nil {
			return err
		}
		result = r
		config.cache.Add(search, result)
	}
	values, err := getExploreValues(result)
	if err != nil {
		return err
	}
	fmt.Println("Exploring: " + values.Name + "...")
	time.Sleep(time.Second * 2)
	fmt.Println("Found Pokemon:")
	for _, poke := range values.PokemonEncounters {
		config.seenPokemon[poke.Pokemon.Name] = values.Name
		fmt.Println(" - " + poke.Pokemon.Name)
	}
	return nil
}
