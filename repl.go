package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartRepl() {
	fmt.Println("Welcome:")
	fmt.Println("type help for commands")
	reader := bufio.NewReader(os.Stdin)
	var globalConfig = Config{"https://pokeapi.co/api/v2/location-area/", nil}
	for {
		fmt.Print("Type > ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\r\n", "", -1)
		command, ok := getCommands()[text]
		if ok {
			err := command.callback(&globalConfig)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknonw Command!")
		}
	}
}
