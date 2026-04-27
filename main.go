/*
Phase 1: The Interstellar Weather Client challenge.

The Objective
Build a Go package that fetches weather data for different planets.
This will teach you how to handle JSON tagging, HTTP timeouts,
and Pointer-based structs (which are used heavily in Terraform to represent "optional" data).
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type WeatherAPI struct {
	Name       *string  `json:"planet_name"`
	Temp       *float64 `json:"temp_celsius"`
	Conditions *string  `json:"conditions"`
	Habitable  *bool    `json:"is_habitable"`
}

type PlanetInfo struct {
	Name                 string
	Temperature          float64
	AtmosphericCondition string
	Habitable            bool
	LastUpdated          time.Time
}

/*
The Task: The "Bridge" Function
Your next goal is to write the function that connects the two worlds. Think of this as the "Ingestion Engine."
This is called Data Normalizer Layer, moving data from "unsafe" pointers into "safe" local variables
*/

func ToInternal(w WeatherAPI) PlanetInfo {

	var planetName string
	if w.Name != nil {
		planetName = *w.Name
	} else {
		planetName = "Mars"
	}

	var planetTemperature float64
	if w.Temp != nil {
		planetTemperature = *w.Temp
	} else {
		planetTemperature = -22.5
	}

	var planetCondition string
	if w.Conditions != nil {
		planetCondition = *w.Conditions
	} else {
		planetCondition = "Dry"
	}

	var planetHabitable bool
	if w.Habitable != nil {
		planetHabitable = *w.Habitable
	} else {
		planetHabitable = false
	}

	planet := PlanetInfo{
		Name:                 planetName,
		Temperature:          planetTemperature,
		AtmosphericCondition: planetCondition,
		Habitable:            planetHabitable,
		LastUpdated:          time.Now(),
	}

	return planet
}

func SaveState(file string, p PlanetInfo) (err error) {
	jsonData, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		fmt.Printf("There was an error enconding the JSON %v\n", err)
		return err
	}
	err = os.WriteFile(file, jsonData, 0644)
	if err != nil {
		fmt.Printf("There was an error writing the JSON into the file(%v) %v\n", file, err)
		return err
	}
	return nil
}

func LoadState(file string) (jsonData PlanetInfo, err error) {
	var p PlanetInfo
	content, err := os.ReadFile(file)
	if err != nil {
		if err == os.ErrNotExist {
			return p, nil //If the error is os.ErrNotExist, treat it as a success case with empty data.
		}
		fmt.Printf("There was an error(%v) reading the contents of the file(%v)\n", file, err)
		return p, err
	}
	err = json.Unmarshal(content, &p)
	if err != nil {
		fmt.Printf("There was an error(%v) decoding the contents of the file(%v)\n", file, err)
		return p, err
	}
	return p, nil
}

func Reconcile(newPlanet PlanetInfo, oldPlanet PlanetInfo, file string) (driftedPlanet PlanetInfo, err error) {

	if newPlanet.Name != oldPlanet.Name {
		fmt.Printf("Drift in Name detected (%v) != (%v)\n", newPlanet.Name, oldPlanet.Name)
	}

	if newPlanet.Temperature != oldPlanet.Temperature {
		fmt.Printf("Drift in Temperature detected (%v) != (%v)\n", newPlanet.Temperature, oldPlanet.Temperature)
	}

	if newPlanet.AtmosphericCondition != oldPlanet.AtmosphericCondition {
		fmt.Printf("Drift in AtmosphericCondition detected (%v) != (%v)\n", newPlanet.AtmosphericCondition, oldPlanet.AtmosphericCondition)
	}

	if newPlanet.Habitable != oldPlanet.Habitable {
		fmt.Printf("Drift in Habitable detected (%v) != (%v)\n", newPlanet.Habitable, oldPlanet.Habitable)
	}

	if newPlanet.LastUpdated.Truncate(time.Minute) != oldPlanet.LastUpdated.Truncate(time.Minute) {
		fmt.Printf("Drift in LastUpdated detected (%v) != (%v)\n", newPlanet.LastUpdated, oldPlanet.LastUpdated)
	}

	err = SaveState(file, newPlanet)

	if err != nil {
		return newPlanet, err
	}

	return newPlanet, nil
}

func main() {
	file := "tfstate.json"
	oldState, err := LoadState(file)

	// How do you turn rawJSON into a WeatherAPI struct?
	// Hint: You'll need "encoding/json" and json.Unmarshal()
	rawJSON := `{"planet_name": "Mars", "temp_celsius": -100.5, "conditions": "Dusty", "is_habitable": false}`

	var jsonBlob = []byte(rawJSON)

	var w WeatherAPI
	err = json.Unmarshal(jsonBlob, &w)
	if err != nil {
		fmt.Println("error: ", err)
	}

	planet := ToInternal(w)
	//fmt.Printf("Planet %v\n", planet)

	driftedPlanet, err := Reconcile(planet, oldState, file)
	if err != nil {
		fmt.Printf("Error detected %v", err)
	}
	fmt.Println(driftedPlanet)

}

