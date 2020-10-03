package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"iaas_sugar/api/sr"
)

func main() {
	seed := flag.String("seed", "", "a file location with some data in JSON form to seed the server content")
	flag.Parse()

	minions := map[string]sr.Minion{}

	if *seed != "" {
		seedData, err := ioutil.ReadFile(*seed)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(seedData, &minions)
		if err != nil {
			log.Fatal(err)
		}
	}

	minionService := sr.NewService("localhost:3001", minions)
	err := minionService.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
