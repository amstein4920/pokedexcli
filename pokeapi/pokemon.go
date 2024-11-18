package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Takes a Pokemon name and returns the fetched Pokemon data and the HTTP and parsing errors
func (c *Client) GetPokemon(name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name

	val, ok := c.cache.Get(url)
	if ok {
		pokemon := Pokemon{}
		err := json.Unmarshal(val, &pokemon)
		if err != nil {
			fmt.Println("Error with Pokemon Cache")
			return Pokemon{}, err
		}
		return pokemon, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error forming request")
		return Pokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println("Error getting response")
		return Pokemon{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading Pokemon")
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	pokemon := Pokemon{}
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		fmt.Println("Error parsing Pokemon")
		return Pokemon{}, err
	}

	c.cache.Add(url, body)

	return pokemon, nil
}
