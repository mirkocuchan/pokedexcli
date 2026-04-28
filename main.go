package main

import (
	"bufio"	
	"os"
	"strings"
	"fmt"
	"log"
	"net/http"
	"io"
	"encoding/json"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}
type config struct{
	Next *string
	Previous *string
}
func commandExit(cfg *config) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error{
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

type LocationArea struct {
	Next *string `json:"next"`
	Previous *string `json:"previous"`
    Results []struct {           //RESOURCE LIST/PAGINATION: Resource Lists/Pagination (group) 
	// Calling any API endpoint without a resource ID or name will return a paginated list of available resources for that API. By default, a list "page" will contain up to 20 resources
        Name string `json:"name"`
        URL  string `json:"url"`
    } `json:"results"`
}

var url = "https://pokeapi.co/api/v2/location-area/"
func commandMap(cfg *config) (error){
	var res *http.Response
	var err error
	if cfg.Next == nil{
		res, err = http.Get(url)
		
	}else{
		res, err = http.Get(*cfg.Next)
		
	}
	if err != nil{
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	var str LocationArea 
	
	err = json.Unmarshal(body, &str)
	if err != nil {
		log.Fatal(err)
	}
	for _, loc := range str.Results {
    	fmt.Println(loc.Name)
	}
	if str.Next != nil {
    	cfg.Next = str.Next
	}
	cfg.Previous = str.Previous	
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
		"map": {
			name:        "map",
			description: "Next or previous 20 location areas",
			callback:    commandMap,
		},
	}
}

func main(){
	reader := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	cfg := &config{}
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

		if err := cmd.callback(cfg); err != nil {
			fmt.Println(err)
		}
		
	}
}

