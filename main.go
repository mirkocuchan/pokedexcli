package main

import "strings"

func cleanInput(text string) []string{
    output := strings.ToLower(text)
    return strings.Fields(output)
}

