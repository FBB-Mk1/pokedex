package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func getLocation(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("You Suck", err)
	}
	body, err := io.ReadAll(res.Body)
	fmt.Println(res.Body)
	fmt.Println(res.StatusCode)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		e := fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return make([]byte, 0), errors.New(e)
	}
	if err != nil {
		return make([]byte, 0), err
	}
	return body, nil
}

func getLocationValues(body []byte) (LocationArea, error) {
	var local LocationArea
	err := json.Unmarshal(body, &local)
	if err != nil {
		return LocationArea{}, err
	}
	return local, nil
}

func getExploreValues(body []byte) (Location, error) {
	var local Location
	err := json.Unmarshal(body, &local)
	if err != nil {
		return Location{}, errors.New("location not found")
	}
	return local, nil
}
