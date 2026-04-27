package main

import (
	"path/filepath"
	"testing"
	"time"
)

// Helper functions that return the mem @
func ptrString(s string) *string {
	return &s
}

func ptrFloat(f float64) *float64 {
	return &f
}

func ptrBool(b bool) *bool {
	return &b
}

func TestToInternal(t *testing.T) {

	//Anonymous struct which has input type of WeatherAPI and expected type of PlanetInfo
	planetWeather := []struct {
		name     string
		input    WeatherAPI
		expected PlanetInfo
	}{
		{
			name: "Happy Path",
			input: WeatherAPI{
				Name:       ptrString("Mars"),
				Temp:       ptrFloat(-65.5),
				Conditions: ptrString("Dusty"),
				Habitable:  ptrBool(false),
			},
			expected: PlanetInfo{
				Name:                 "Mars",
				Temperature:          -65.5,
				AtmosphericCondition: "Dusty",
				Habitable:            false,
				LastUpdated:          time.Now(),
			},
		},
		{
			name: "Handle missing name",
			input: WeatherAPI{
				Temp:       ptrFloat(-65.5),
				Conditions: ptrString("Dusty"),
				Habitable:  ptrBool(false),
			},
			expected: PlanetInfo{
				Name:                 "Mars",
				Temperature:          -65.5,
				AtmosphericCondition: "Dusty",
				Habitable:            false,
				LastUpdated:          time.Now(),
			},
		},
		{
			name: "Handle missing temperature",
			input: WeatherAPI{
				Name:       ptrString("Mars"),
				Conditions: ptrString("Dusty"),
				Habitable:  ptrBool(false),
			},
			expected: PlanetInfo{
				Name:                 "Mars",
				Temperature:          -22.5,
				AtmosphericCondition: "Dusty",
				Habitable:            false,
				LastUpdated:          time.Now(),
			},
		},
		{
			name: "Handle missing atmospheric condition",
			input: WeatherAPI{
				Name:      ptrString("Mars"),
				Temp:      ptrFloat(-65.5),
				Habitable: ptrBool(false),
			},
			expected: PlanetInfo{
				Name:                 "Mars",
				Temperature:          -65.5,
				AtmosphericCondition: "Dry",
				Habitable:            false,
				LastUpdated:          time.Now(),
			},
		},
		{
			name: "Handle missing habitability",
			input: WeatherAPI{
				Name:       ptrString("Mars"),
				Temp:       ptrFloat(-65.5),
				Conditions: ptrString("Dusty"),
			},
			expected: PlanetInfo{
				Name:                 "Mars",
				Temperature:          -65.5,
				AtmosphericCondition: "Dusty",
				Habitable:            false,
				LastUpdated:          time.Now(),
			},
		},
	}

	for _, tc := range planetWeather {
		t.Run(tc.name, func(t *testing.T) {
			got := ToInternal(tc.input)
			if got.Name != tc.expected.Name {
				t.Errorf("Got(%v) but want(%v)", got.Name, tc.expected.Name)
			}
			if got.Temperature != tc.expected.Temperature {
				t.Errorf("Got(%v) but want(%v)", got.Temperature, tc.expected.Temperature)
			}
			if got.AtmosphericCondition != tc.expected.AtmosphericCondition {
				t.Errorf("Got(%v) but want(%v)", got.AtmosphericCondition, tc.expected.AtmosphericCondition)
			}
			if got.Habitable != tc.expected.Habitable {
				t.Errorf("Got(%v) but want(%v)", got.Habitable, tc.expected.Habitable)
			}
			if got.Habitable != tc.expected.Habitable {
				t.Errorf("Got(%v) but want(%v)", got.Habitable, tc.expected.Habitable)
			}
			if got.LastUpdated.IsZero() {
				t.Error("Time was never set")
			}
		})
	}
}

func TestStateFile(t *testing.T) {
	p := PlanetInfo{
		Name:                 "Jupiter",
		Temperature:          -100.10,
		AtmosphericCondition: "Rocky",
		Habitable:            false,
		LastUpdated:          time.Now(),
	}

	tempFile := filepath.Join(t.TempDir(), "test_state.json")
	SaveState(tempFile, p)
	jsonData, err := LoadState(tempFile)
	if err != nil {
		t.Fatal(err)
	}
	if jsonData.Name != p.Name {
		t.Errorf("Got(%v) but want(%v)", jsonData.Name, p.Name)
	}
	if jsonData.Temperature != p.Temperature {
		t.Errorf("Got(%v) but want(%v)", jsonData.Temperature, p.Temperature)
	}
	if jsonData.AtmosphericCondition != p.AtmosphericCondition {
		t.Errorf("Got(%v) but want(%v)", jsonData.AtmosphericCondition, p.AtmosphericCondition)
	}
	if jsonData.Habitable != p.Habitable {
		t.Errorf("Got(%v) but want(%v)", jsonData.Habitable, p.Habitable)
	}
	if jsonData.LastUpdated.Truncate(time.Second) != p.LastUpdated.Truncate(time.Second) {
		t.Errorf("Got(%v) but want(%v)", jsonData.LastUpdated, p.LastUpdated)
	}
}
