package openweathermap

import (
	"testing"
)

func TestGetWeatherDataFromApi(t *testing.T) {
	city := "Berlin"
	_, err := GetWeatherDataFromApi(city)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	city = "HHHMoscow"
	_, err = GetWeatherDataFromApi(city)

	if err.Error()!="no such city" {
		t.Errorf("Expected an error, got %v", err)
	}

}
