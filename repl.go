package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fbb-mk1/pokedex/internal/pkcache"
)

func StartRepl(cache pkcache.Cache) {
	fmt.Println("Welcome:")
	fmt.Println("type help for commands")
	reader := bufio.NewReader(os.Stdin)
	var globalConfig = Config{cache,
		"https://pokeapi.co/api/v2/location-area/",
		nil,
		map[string]PokedexEntry{},
	}

	for {
		fmt.Print("Pokedex > ")
		input, _ := reader.ReadString('\n')
		text := strings.Split(strings.Replace(input, "\r\n", "", -1), " ")
		command, ok := getCommands()[text[0]]
		if ok {
			arg := " "
			if len(text) > 1 {
				arg = text[1]
			}
			err := command.callback(&globalConfig, arg)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknonw Command!")
		}
	}
}
