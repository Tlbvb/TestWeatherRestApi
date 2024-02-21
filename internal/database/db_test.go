package database

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/tlbvb/weatherestapi/openweathermap"
)

func PrepareTestApis(db *pgx.Conn) {
	qs := []string{
		`DROP TABLE IF EXISTS weath;`,

		`CREATE TABLE weath (
			cityname VARCHAR(100) PRIMARY KEY NOT NULL,
			temp FLOAT NOT NULL,
			feels_like FLOAT NOT NULL,
			temp_min FLOAT NOT NULL,
			temp_max FLOAT NOT NULL,
			humidity INT NOT NULL,
			pressure INT NOT NULL
		);`,

		` INSERT INTO weath (cityname, temp, feels_like, temp_min, temp_max, humidity, pressure) VALUES
		('New York', 25.5, 28.0, 22.0, 30.0, 70, 1012),
		('London', 18.2, 20.5, 15.0, 22.0, 65, 1018),
		('Tokyo', 30.8, 32.5, 28.0, 34.0, 80, 1005),
		('TestCity', 25.0, 27.0, 20.0, 30.0, 70, 1010);`,
		
	}

	for _, q := range qs {
		_, err := db.Exec(context.Background(),q)
		if err != nil {
			panic(err)
		}
	}
}

func TestConnect(t *testing.T) {
	connStr := "postgres://testUser:test@localhost:5432/weather?sslmode=disable"
	conn, err := Connect(connStr)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	conn.Close(context.Background())

	connStr = "postgres://testUser:test@localhost:5432/city?sslmode=disable"
	conn, err = Connect(connStr)
	if err != nil {
		if err.Error()!="connection error"{
			t.Fatalf("Not the right error")
			return
		}
	}
	if conn!=nil{
		conn.Close(context.Background())
	}
}

func TestGetWeatherData(t *testing.T) {
	connStr := "postgres://testUser:test@localhost:5432/weather?sslmode=disable"
	conn, err := Connect(connStr)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer conn.Close(context.Background())
	PrepareTestApis(conn)
	cityName := "TestCity"
	expectedData := &openweathermap.WeatherData{
		Temp:      25.0,
		FeelsLike: 27.0,
		MinTemp:   20.0,
		MaxTemp:   30.0,
		Pressure:  1010,
		Humidity:  70,
	}

	var wData openweathermap.WeatherData
	err = GetWeatherData(conn, cityName, &wData)
	if err != nil {
		t.Fatalf("GetWeatherData failed: %v", err)
	}

	if wData.Temp != expectedData.Temp ||
		wData.FeelsLike != expectedData.FeelsLike ||
		wData.MinTemp != expectedData.MinTemp ||
		wData.MaxTemp != expectedData.MaxTemp ||
		wData.Pressure != expectedData.Pressure ||
		wData.Humidity != expectedData.Humidity {
		t.Fatalf("Unexpected result. Expected: %+v, Actual: %+v", expectedData, wData)
	}
}

func TestUpdateWeatherData(t *testing.T) {
	connStr := "postgres://testUser:test@localhost:5432/weather?sslmode=disable"
	conn, err := Connect(connStr)
	if err != nil {
		t.Fatalf("Connect failed: %v", err)
	}
	defer conn.Close(context.Background())

	PrepareTestApis(conn)
	cityName := "London"

	_, err = UpdateWeatherData(conn, cityName)
	if err != nil {
		t.Fatalf("UpdateWeatherData failed: %v", err)
	}

	cityName = "Astana"

	_, err = UpdateWeatherData(conn, cityName)
	if err != nil {
		t.Fatalf("UpdateWeatherData failed: %v", err)
	}
}


