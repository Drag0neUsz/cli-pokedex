package main

import pokecache "cli-pokedex/internal"

type cliCommand struct {
	name           string
	description    string
	callback       func(*config) error
	callbackConfig *config
}

type config struct {
	nextUrl     string
	previousUrl string
	pokecache   *pokecache.Cache
}

type locationAreaResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
