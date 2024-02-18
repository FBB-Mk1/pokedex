package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func getLocation(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("You Suck", err)
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		e := fmt.Sprintf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		log.Fatal(e)
	}
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func getLocationValues(body []byte) LocationArea {
	var local LocationArea
	err := json.Unmarshal(body, &local)
	if err != nil {
		log.Fatal(err)
	}
	return local
}
