package main

import (
	pokecache "cli-pokedex/internal"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	LocationAreaApiUrl = "https://pokeapi.co/api/v2/location-area"
	PokemonApiUrl      = "https://pokeapi.co/api/v2/pokemon"
	CacheDuration      = 7 * time.Second
)

func getCommands() map[string]cliCommand {
	mapConfig := &config{nextUrl: LocationAreaApiUrl, previousUrl: "", pokecache: pokecache.NewCache(CacheDuration)}
	pokedexConfig := &config{pokecache: pokecache.NewCache(CacheDuration), pokedex: make(map[string]pokemonType)}
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
			callbackConfig: pokedexConfig,
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
		"catch": {
			name:           "catch",
			description:    "Catch a Pokemon passed as an argument <pokemon-name>",
			callback:       commandCatch,
			callbackConfig: pokedexConfig,
		},
		"inspect": {
			name:           "inspect",
			description:    "Inspect a Pokemon passed as an argument <pokemon-name>",
			callback:       commandInspect,
			callbackConfig: pokedexConfig,
		},
		"pokedex": {
			name:           "pokedex",
			description:    "Display the Pokedex (or wipe it using 'clear' as an argument)",
			callback:       commandPokedex,
			callbackConfig: pokedexConfig,
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	err := savePokedex(cfg.pokedex)
	if err != nil {
		fmt.Printf("error saving pokedex: %v\n", err)
	}
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
		fmt.Printf("  - %v\n", result.Name)
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
		fmt.Printf("  - %v\n", result.Name)
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
		fmt.Printf("  - %v\n", result.Pokemon.Name)
	}
	resBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %v", err)
	}
	cfg.pokecache.Add(areaName, resBytes)
	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no pokemon name passed")
	}
	pokemonUrl := fmt.Sprintf("%s/%s", PokemonApiUrl, args[0])
	var data pokemonResponse
	resBytes, ok := cfg.pokecache.Get(pokemonUrl)

	if !ok {
		res, err := http.Get(pokemonUrl)
		if err != nil {
			return fmt.Errorf("error fetching pokemon: %v", err)
		}
		defer res.Body.Close()

		if res.StatusCode != 200 {
			return fmt.Errorf("error fetching pokemon (code: %v)", res.StatusCode)
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

	pokemonData := pokemonType{
		Name:   data.Name,
		Weight: data.Weight,
		Stats:  data.Stats,
		Types:  data.Types,
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonData.Name)
	randomNumber := rand.Float64()
	chance := math.Exp(-0.0069315 * float64(data.BaseExperience)) // 50% at 100 base experience maybe i'll change this later (36 min to 600 max afaik)
	if randomNumber <= chance {
		fmt.Printf("%v was caught!\n", pokemonData.Name)
		cfg.pokedex[pokemonData.Name] = pokemonData
	} else {
		fmt.Printf("%v escaped!\n", pokemonData.Name)
	}
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no pokemon name passed")
	}
	pokemonName := args[0]
	pokemonData, ok := cfg.pokedex[pokemonName]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}
	fmt.Printf("Name: %v\n", pokemonData.Name)
	fmt.Printf("Weight: %v\n", pokemonData.Weight)
	fmt.Printf("Stats: \n")
	for _, stat := range pokemonData.Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types: \n")
	for _, typeData := range pokemonData.Types {
		fmt.Printf("  -%v\n", typeData.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *config, args ...string) error {
	fmt.Printf("Pokedex: \n")
	if len(args) > 0 {
		if args[0] == "clear" {
			cfg.pokedex = make(map[string]pokemonType)
			fmt.Printf("  - Pokedex cleared\n")
			return nil
		}
	}
	if len(cfg.pokedex) == 0 {
		fmt.Printf("  - No Pokemon caught yet\n")
		return nil
	}
	for pokemonName, _ := range cfg.pokedex {
		fmt.Printf("  - %v\n", pokemonName)
	}
	return nil
}
