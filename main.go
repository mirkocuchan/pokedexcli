package main

import (
	"bufio"	
	"os"
	"strings"
	"fmt"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandExit() error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error{
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func cleanInput(text string) []string{
    output := strings.ToLower(text)
    return strings.Fields(output)
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
	}
}

func main(){
	reader := bufio.NewScanner(os.Stdin)
	commands := getCommands()

	for{
		fmt.Print("Pokedex > ")
		if !reader.Scan() {
    		break
		}
		
		words := cleanInput(reader.Text())
			
		if len(words) == 0 {
			continue 
		}
		
		cmd, ok := commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := cmd.callback(); err != nil {
			fmt.Println(err)
		}
		
	}
}

