package main

import (
	pokecache "cli-pokedex/internal"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	LocationAreaApiUrl = "https://pokeapi.co/api/v2/location-area"
	CacheDuration      = 7 * time.Second
)

func getCommands() map[string]cliCommand {
	mapConfig := &config{nextUrl: LocationAreaApiUrl, previousUrl: "", pokecache: pokecache.NewCache(CacheDuration)}
	return map[string]cliCommand{
		"help": {
			name:           "help",
			description:    "Displays a help message",
			callback:       commandHelp,
			callbackConfig: &config{},
		},
		"exit": {
			name:           "exit",
			description:    "Exit the Pokedex",
			callback:       commandExit,
			callbackConfig: &config{},
		},
		"map": {
			name:           "map",
			description:    "Display next page of locations ",
			callback:       commandMap,
			callbackConfig: mapConfig,
		},
		"mapb": {
			name:           "map-back",
			description:    "Display previous page of locations (if available)",
			callback:       commandMapb,
			callbackConfig: mapConfig,
		},
		"explore": {
			name:           "explore",
			description:    "Explore the area passed as an argument <area-name>",
			callback:       commandExplore,
			callbackConfig: &config{pokecache: pokecache.NewCache(CacheDuration)},
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for key, obj := range getCommands() {
		fmt.Printf("%v: %v\n", key, obj.description)
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
	var data locationAreaResponse
	resBytes, ok := cfg.pokecache.Get(cfg.nextUrl)

	if !ok {
		res, err := http.Get(cfg.nextUrl)
		if err != nil {
			return fmt.Errorf("error fetching next url: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return fmt.Errorf("error fetching next url (code: %v)", res.StatusCode)
		}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return fmt.Errorf("error reading body: %v", err)
		}
	} else {
		err := json.Unmarshal(resBytes, &data)
		if err != nil {
			return fmt.Errorf("error reading body: %v", err)
		}
	}

	for _, result := range data.Results {
		fmt.Printf("%v\n", result.Name)
	}
	resBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}
	cfg.pokecache.Add(cfg.nextUrl, resBytes)
	cfg.nextUrl = data.Next
	cfg.previousUrl = data.Previous

	// fmt.Printf("Next: %v\n", cfg.nextUrl)
	// fmt.Printf("Previous: %v\n", cfg.previousUrl)
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.previousUrl == "" {
		fmt.Printf("you're on the first page\n")
		return nil
	}
	var data locationAreaResponse
	resBytes, ok := cfg.pokecache.Get(cfg.previousUrl)
	if !ok {
		res, err := http.Get(cfg.previousUrl)
		if err != nil {
			return fmt.Errorf("error fetching previous url: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return fmt.Errorf("error fetching previous url (code: %v)", res.StatusCode)
		}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return fmt.Errorf("error reading body: %v", err)
		}
	} else {
		err := json.Unmarshal(resBytes, &data)
		if err != nil {
			return fmt.Errorf("error reading body: %v", err)
		}
	}

	for _, result := range data.Results {
		fmt.Printf("%v\n", result.Name)
	}
	resBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}
	cfg.pokecache.Add(cfg.previousUrl, resBytes)
	cfg.nextUrl = data.Next
	cfg.previousUrl = data.Previous
	return nil

}

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no area name passed")
	}
	var data pokemonEncounterResponse
	areaName := args[0]
	resBytes, ok := cfg.pokecache.Get(areaName)
	if !ok {
		res, err := http.Get(fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", areaName))
		if err != nil {
			return fmt.Errorf("error fetching area: %v", err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return fmt.Errorf("error fetching area (code: %v)", res.StatusCode)
		}
		err = json.NewDecoder(res.Body).Decode(&data)
		if err != nil {
			return fmt.Errorf("error reading body: %v", err)
		}
	} else {
		err := json.Unmarshal(resBytes, &data)
		if err != nil {
			return fmt.Errorf("error reading body: %v", err)
		}
	}
	fmt.Printf("Exploring %s...\n", areaName)
	fmt.Printf("Found Pokemon:\n")
	if len(data.PokemonEncounters) == 0 {
		fmt.Printf("No Pokemon found\n")
		return nil
	}
	for _, result := range data.PokemonEncounters {
		fmt.Printf("%v\n", result.Pokemon.Name)
	}
	resBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}
	cfg.pokecache.Add(areaName, resBytes)
	return nil
}
