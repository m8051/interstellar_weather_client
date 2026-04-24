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

func main() {
	// How do you turn rawJSON into a WeatherAPI struct?
	// Hint: You'll need "encoding/json" and json.Unmarshal()
	rawJSON := `{"planet_name": "Mars", "temp_celsius": -65.5, "conditions": "Dusty", "is_habitable": false}`

	var jsonBlob = []byte(rawJSON)

	var w WeatherAPI
	err := json.Unmarshal(jsonBlob, &w)
	if err != nil {
		fmt.Println("error: ", err)
	}

	planet := ToInternal(w)
	fmt.Printf("Planet %v\n", planet)
}
