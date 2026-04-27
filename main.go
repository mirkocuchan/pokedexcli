package main

import (
	"bufio"	
	"os"
	"strings"
	"fmt"
)

func cleanInput(text string) []string{
    output := strings.ToLower(text)
    return strings.Fields(output)
}

func main(){
	reader := bufio.NewScanner(os.Stdin)
	for{
		fmt.Print("Pokedex > ")
		reader.Scan()
		
		words := cleanInput(reader.Text())
			
		if len(words) == 0 {
			continue 
		}

		fmt.Println("Your command was:", words[0])
	}
}