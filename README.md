# Overview

This is a simple cli project written in Go that acts as a pokedex. It uses the PokeAPI to get the data for the pokemon.
I've also implemented a simple cache to store the data in memory.

Current commands:

- help: Display the help menu
- map: Display 20 (next) locations
- mapb: Display the previous 20 locations (if applicable)
- explore `<area-name>` Explore the area `<area-name>`
- catch `<pokemon-name>`: Catch a pokemon `<pokemon-name>`
- inspect `<pokemon-name>`: Inspect a pokemon `<pokemon-name>`
- pokedex: Display the pokedex
- exit: Exit the program

## Usage
### Prerequisites

* Go 1.20+ installed on your system.

### Running the Application

To start the Pokedex CLI, run the following command in the root directory of the project:

```bash
go run .
```

### Example Session
```bash
Pokedex > map
- canalave-city-area
- eterna-city-area
- pastoria-city-area
...

Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found pokemon:
- tentacool
- tentacruel
- magikarp

Pokedex > catch magikarp
Throwing a Pokeball at magikarp...
magikarp escaped!

Pokedex > catch magikarp
Throwing a Pokeball at magikarp...
magikarp was caught!

Pokedex > inspect magikarp
Name: magikarp
Weight: 100
Stats: 
  -hp: 20
  -attack: 10
  -defense: 55
  -special-attack: 15
  -special-defense: 20
  -speed: 80
Types: 
  -water

Pokedex > pokedex
Pokedex: 
  - magikarp

Pokedex > exit
Closing the Pokedex... Goodbye!
```