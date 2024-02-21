package openweathermap

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)
type Coordinates struct{
	Latitude float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

type WeatherData struct{
	Temp float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	MinTemp float64 `json:"temp_min"`
	MaxTemp float64 `json:"temp_max"`
	Pressure int `json:"pressure"`
	Humidity int `json:"humidity"`
	City string `json:"-"`
}

func GetWeatherDataFromApi(city string) (*WeatherData,error){
	apiKey:="e46de9febe313e42050a8ab329d89835"
	fmt.Println("cc",city)
	urlCord:=fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s",city,apiKey)
	resp,err:=http.Get(urlCord)
	if err!=nil{
		return nil,errors.New("no such city")
	}
	fmt.Println("r",resp)
	fmt.Println("e",err)
	defer resp.Body.Close()
	var cord []Coordinates
	err = json.NewDecoder(resp.Body).Decode(&cord)
	if err != nil {
		return nil,err
	}
	if len(cord)==0{
		return nil,errors.New("no such city")
	}
	fmt.Println(cord)
	url:=fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric",cord[0].Latitude,cord[0].Longitude,apiKey)
	resp,err=http.Get(url)
	if err!=nil{
		return nil,err
	}
	fmt.Println("resp",resp)
	var x map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&x)
	if err != nil {
		return nil,err
	}

	mainB,err:=json.Marshal(x["main"])
	if err != nil {
		return nil,err
	}
	wData:=WeatherData{}
	json.Unmarshal(mainB,&wData)
	return &wData,nil
}