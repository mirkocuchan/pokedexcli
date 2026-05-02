package main

import (
	"bufio"	
	"os"
	"strings"
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"errors"
	"github.com/mirkocuchan/pokedexcli/internal/pokecache"
	"time"
	"math/rand"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}
type config struct{
	Next *string
	Previous *string
	currentLocation        string
    currentLocationPokemon map[string]bool
	pokedex                *Pokedex
}
func commandExit(cfg *config, args []string) error{
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args []string) error{
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
var cache = pokecache.NewCache(5 * time.Second)

func commandMap(cfg *config, args []string) (error){
	currentURL := url

	if cfg.Next != nil {
    	currentURL = *cfg.Next
	}
	
	var body []byte
    data, ok := cache.Get(currentURL)

    if ok {
        body = data
    } else {
        res, err := http.Get(currentURL)
        if err != nil {
            return err
        }
        defer res.Body.Close()

        body, err = io.ReadAll(res.Body)
        if err != nil {
            return err
        }

        cache.Add(currentURL, body)
    }

    var str LocationArea
    err := json.Unmarshal(body, &str)
    if err != nil {
        return err
    }

    for _, loc := range str.Results {
        fmt.Println(loc.Name)
    }

    cfg.Next = str.Next
    cfg.Previous = str.Previous

    return nil
}

func commandMapb(cfg *config, args []string) (error){

	currentURL := url
	if cfg.Previous == nil{
		return errors.New("you're on the first page")
	}
	currentURL = *cfg.Previous
	
	
	var err error
	var body []byte
    data, ok := cache.Get(currentURL)

	if ok{
		body = data
	}else{
		res, err := http.Get(currentURL)
        if err != nil {
            return err
        }
        defer res.Body.Close()

        body, err = io.ReadAll(res.Body)
        if err != nil {
            return err
        }

        cache.Add(currentURL, body)
	}	
	
	var str LocationArea 
	
	err = json.Unmarshal(body, &str)
	if err != nil {
		return err
	}
	for _, loc := range str.Results {
    	fmt.Println(loc.Name)
	}
    cfg.Next = str.Next

	cfg.Previous = str.Previous	
	return nil 
}

type LocationAreaDetail struct {
	Name              string `json:"name"`
    PokemonEncounters []struct {
        Pokemon struct {
            Name string `json:"name"`
            URL  string `json:"url"`
        } `json:"pokemon"`
    } `json:"pokemon_encounters"`
}

func commandExplore(cfg *config, args []string) error{
	if len(args) != 1 {
    	return errors.New("you must provide a location name")
	}
	name := args[0]
	url := "https://pokeapi.co/api/v2/location-area/" + name + "/"
	var body []byte
	var err error

    data, ok := cache.Get(url)
	if ok{
		body = data
	}else{
		res, err := http.Get(url)
		if err != nil {
            return err
        }

        defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
        if err != nil {
            return err
        }

        cache.Add(url, body)
	}

	var str LocationAreaDetail 
	
	err = json.Unmarshal(body, &str)
	if err != nil {
		return err
	}
	cfg.currentLocationPokemon = map[string]bool{}
	cfg.currentLocation = str.Name

	fmt.Printf("Exploring %s...\n", name)
	fmt.Println("Found Pokemon: ")
	for _, pokemon := range str.PokemonEncounters {
    	fmt.Println(pokemon.Pokemon.Name)
		cfg.currentLocationPokemon[pokemon.Pokemon.Name] = true
	}
    
	return nil 
}

func commandCatch(cfg *config, args []string) error{
	if len(args) != 1 {
    	return errors.New("you must provide a pokemon name")
	}
	if cfg.currentLocation == "" {
    	return errors.New("you must explore a location first")
	}
	name := args[0]
	if !cfg.currentLocationPokemon[name] {
    	return fmt.Errorf("%s is not in %s", name, cfg.currentLocation)
	}
	
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	url := "https://pokeapi.co/api/v2/pokemon/" + name + "/"
	
	var body []byte
	var err error
	data, ok := cache.Get(url)
	if ok{
    	body = data
	}else{
    	res, err := http.Get(url)
		if err != nil {
            return err
        }

        defer res.Body.Close()

		body, err = io.ReadAll(res.Body)
        if err != nil {
            return err
        }

    	cache.Add(url, body)
	}
	
	var str Pokemon
	err = json.Unmarshal(body, &str)
	if err != nil {
		return err
	}

	chance := calcularChance(str.BaseExperience)
	randomNumber := rand.Float64()
	if chance > randomNumber{
		fmt.Printf("%s was caught!\n", name)
		fmt.Printf("You may now inspect it with the inspect command.")
		cfg.pokedex.Add(str)
	}else{
		fmt.Printf("%s escaped!\n", name)
	}    
	return nil 
}

func commandInspect(cfg *config, args []string) error{
	name := args[0]

	data, ok := cfg.pokedex.Get(name)
	if !ok{
		return errors.New("you don't have it lil bro")
	}
	fmt.Println("Name:", data.Name)
    fmt.Println("Height:", data.Height)
    fmt.Println("Weight:", data.Weight)
	fmt.Println("Stats:")
	for _, stat := range data.Stats{
		fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typecito := range data.Types {
    	fmt.Printf("  - %s\n", typecito.Type.Name)
	}
	return nil
}
func calcularChance(baseExperience int) float64{
	return 1.0 / float64(baseExperience + 1)
}

func commandPokedex(cfg *config, args []string) error{
	fmt.Printf("Your Pokedex:\n")
	for _, pokemon := range cfg.pokedex.Pokemons{
		fmt.Printf("  - %s\n", pokemon.Name)
	}
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
			description: "Next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "List of Pokemon in the area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Get the details",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "See your pokemons",
			callback:    commandPokedex,
		},
	}
}

func main(){
	reader := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	cfg := &config{
		currentLocationPokemon: map[string]bool{},
		pokedex:                NewPokedex(),
	}
	args := []string{}
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
		if len(words) > 1 {
    		args = words[1:]
		} else {
    		args = []string{}
		}
		
		if err := cmd.callback(cfg, args); err != nil {
				fmt.Println(err)
		}			
	}
}

