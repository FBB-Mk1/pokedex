package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fbb-mk1/pokedex/internal/pkcache"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	cache    pkcache.Cache
	next     string
	previous *string
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
			description: "Show the previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func commandMap(config *Config) error {
	local, ok := config.cache.Get(config.next)
	if !ok {
		local = getLocation(config.next)
	}
	l := getLocationValues(local)
	config.next = l.Next
	config.previous = l.Previous
	for _, loc := range l.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandMapb(config *Config) error {
	if config.previous == nil {
		return errors.New("at start of list")
	}
	local, ok := config.cache.Get(*config.previous)
	if !ok {
		local = getLocation(*config.previous)
	}
	l := getLocationValues(local)
	config.next = l.Next
	config.previous = l.Previous
	for _, loc := range l.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func commandHelp(*Config) error {
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

func commandExit(*Config) error {
	fmt.Println("Bye Bye")
	os.Exit(0)
	return nil
}
