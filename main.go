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

func ToInternal(w []WeatherAPI) map[string]PlanetInfo {

	planetMap := make(map[string]PlanetInfo)
	for _, j := range w {
		var planetName string
		if j.Name != nil {
			planetName = *j.Name
		} else {
			planetName = "Mars"
		}

		var planetTemperature float64
		if j.Temp != nil {
			planetTemperature = *j.Temp
		} else {
			planetTemperature = -22.5
		}

		var planetCondition string
		if j.Conditions != nil {
			planetCondition = *j.Conditions
		} else {
			planetCondition = "Dry"
		}

		var planetHabitable bool
		if j.Habitable != nil {
			planetHabitable = *j.Habitable
		} else {
			planetHabitable = false
		}

		planet := PlanetInfo{
			Temperature:          planetTemperature,
			AtmosphericCondition: planetCondition,
			Habitable:            planetHabitable,
			LastUpdated:          time.Now(),
		}
		planetMap[planetName] = planet
	}

	return planetMap
}

func SaveState(file string, p map[string]PlanetInfo) (err error) {
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

func LoadState(file string) (jsonData map[string]PlanetInfo, err error) {
	var p map[string]PlanetInfo
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

func Reconcile(newPlanet map[string]PlanetInfo, oldPlanet map[string]PlanetInfo) (plans []string, p map[string]PlanetInfo) {
	var changes []string

	for name, newData := range newPlanet {
		oldData, exists := oldPlanet[name]
		if !exists {
			changes = append(changes, fmt.Sprintf("Initial state creation for %s", name))
		} else {

			if newData.Temperature != oldData.Temperature {
				changes = append(changes, fmt.Sprintf("Temperature drift: %f -> %f", oldData.Temperature, newData.Temperature))
			}
			if newData.AtmosphericCondition != oldData.AtmosphericCondition {
				changes = append(changes, fmt.Sprintf("AtmosphericCondition drift: %s -> %s", oldData.AtmosphericCondition, newData.AtmosphericCondition))
			}
			if newData.Habitable != oldData.Habitable {
				changes = append(changes, fmt.Sprintf("Habitable drift: %v -> %v", oldData.Habitable, newData.Habitable))
			}
		}
	}

	// We skip Time for the test to avoid "flaky" results based on milliseconds
	return changes, newPlanet
}

func main() {
	file := "tfstate.json"
	oldState, err := LoadState(file)

	rawJSON := `[
		{ "planet_name": "Mars", "temp_celsius": -65.5, "conditions": "Dusty","is_habitable": false },
		{ "planet_name": "Venus", "temp_celsius": 464.0, "conditions": "Acid Rain", "is_habitable": false },
		{ "planet_name": "Earth", "temp_celsius": 15.0, "conditions": "Cloudy", "is_habitable": true }
	]`

	var jsonBlob = []byte(rawJSON)

	var w []WeatherAPI
	err = json.Unmarshal(jsonBlob, &w)
	if err != nil {
		fmt.Println("error: ", err)
	}

	planet := ToInternal(w)
	fmt.Printf("Planet %v\n", planet)

	plan, updatedPlanet := Reconcile(planet, oldState)

	for _, change := range plan {
		fmt.Println("PLAN:", change)
	}

	if len(plan) > 0 {
		err := SaveState(file, updatedPlanet)
		if err != nil {
			fmt.Printf("Error saving the content in the file (%v):(%v)\n", file, err)
		}
	}

}
