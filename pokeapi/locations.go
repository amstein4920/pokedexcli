package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokeLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) ListLocations(nextUrl *string) (PokeLocations, error) {
	url := baseURL + "/location-area"
	if nextUrl != nil {
		url = *nextUrl
	}

	val, ok := c.cache.Get(url)
	if ok {
		locations := PokeLocations{}
		err := json.Unmarshal(val, &locations)
		if err != nil {
			fmt.Println("Error with cached locations JSON")
			return PokeLocations{}, err
		}

		return locations, nil
	}

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error loading locations")
		return PokeLocations{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading locations")
		return PokeLocations{}, err
	}
	defer res.Body.Close()

	locations := PokeLocations{}
	err = json.Unmarshal(body, &locations)
	if err != nil {
		fmt.Println("Error with locations JSON")
		return PokeLocations{}, err
	}
	return locations, nil
}
