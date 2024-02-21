package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/tlbvb/weatherestapi/openweathermap"
)


func Connect(connstr string) (*pgx.Conn,error){
	conn,err:=pgx.Connect(context.Background(),connstr)
	fmt.Println("conn",conn)
	if err!=nil{
		return nil,errors.New("сonnection error")
	}
	return conn,nil
}

func GetWeatherData(conn *pgx.Conn,cityName string,wdata *openweathermap.WeatherData) (error){
	row := conn.QueryRow(context.Background(), "SELECT temp, feels_like, temp_min, temp_max, pressure, humidity FROM weath WHERE cityname=$1	;", cityName)
	fmt.Println("row",row)

	if err := row.Scan(&wdata.Temp, &wdata.FeelsLike, &wdata.MinTemp, &wdata.MaxTemp, &wdata.Pressure, &wdata.Humidity); err != nil {
		return err
    }
	fmt.Printf("Values from Scan: %v %v %v %v %v %v\n", wdata.Temp, wdata.FeelsLike, wdata.MinTemp, wdata.MaxTemp, wdata.Pressure, wdata.Humidity)
	fmt.Println("wdata",wdata)
	return nil
}


func UpdateWeatherData(conn *pgx.Conn, cityName string) (openweathermap.WeatherData,error) {
    wData,err:= openweathermap.GetWeatherDataFromApi(cityName)
	fmt.Printf("%T %s",cityName,cityName)
	if err!=nil{
		fmt.Println("eee",err)
		return *wData, errors.New("failed to fetch weather data from API")
	}
    
    query := `
        UPDATE weath
        SET temp=$1, feels_like=$2, temp_min=$3, temp_max=$4, humidity=$5, pressure=$6 
        WHERE cityname=$7
    `

    commandTag, err := conn.Exec(
        context.Background(),
        query,
        wData.Temp, wData.FeelsLike, wData.MinTemp, wData.MaxTemp, wData.Humidity, wData.Pressure, cityName,
    )
	
    if err != nil {
        return *wData,err
    }
	if commandTag.RowsAffected()==0{
		query = `
    INSERT INTO weather 
    (cityname, temp, feels_like, temp_min, temp_max, humidity, pressure) 
    VALUES ($1, $2, $3, $4, $5, $6, $7)
`
		_,_=conn.Exec(context.Background(),query,cityName,wData.Temp, wData.FeelsLike, wData.MinTemp, wData.MaxTemp, wData.Humidity, wData.Pressure)
	}
    return *wData,nil
}
