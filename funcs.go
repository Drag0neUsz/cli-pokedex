package main

import (
	"encoding/json"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	return words
}

func savePokedex(pokedex map[string]pokemonType) error {
	jsonData, err := json.Marshal(pokedex)
	if err != nil {
		return err
	}
	return os.WriteFile("pokedex.json", jsonData, 0644)
}

func loadPokedex() (map[string]pokemonType, error) {
	jsonData, err := os.ReadFile("pokedex.json")
	if err != nil {
		return nil, err
	}
	var pokedex map[string]pokemonType
	err = json.Unmarshal(jsonData, &pokedex)
	if err != nil {
		return nil, err
	}
	return pokedex, nil
}
