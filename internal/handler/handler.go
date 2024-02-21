package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/tlbvb/weatherestapi/internal/database"
	"github.com/tlbvb/weatherestapi/openweathermap"
)


type WeatherHandler struct{
	Conn *pgx.Conn
}


func (wh *WeatherHandler) ServeHTTP(w http.ResponseWriter,r *http.Request){
	if r.URL.Path=="/weather"{
		if r.Method=="GET"{
			cityName:=r.FormValue("city")
			if cityName==""{
				w.WriteHeader(http.StatusBadRequest)
				m,_:=json.Marshal("didn't mention the city")
				w.Write(m)
				return
			}
			fmt.Println(cityName)
			Wdata:=&openweathermap.WeatherData{}
			err:=database.GetWeatherData(wh.Conn,cityName,Wdata )
			if err!=nil{
				w.WriteHeader(http.StatusBadRequest)
				fmt.Println("no such")
				m,_:=json.Marshal("No such city in table")
				w.Write(m)
				return
			}
			fmt.Println("wdata",Wdata)
			Wdata.City=cityName
			retB,err:=json.Marshal(Wdata)
			if err!=nil{
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(retB)
		}else if r.Method=="PUT"{
			cityName:=r.FormValue("city")
			fmt.Println(wh.Conn)
			fmt.Println("//",cityName)
			if cityName==""{
				w.WriteHeader(http.StatusBadRequest)
				m,_:=json.Marshal("didn't mention the city")
				w.Write(m)
				return
			}
			wdat,err:=database.UpdateWeatherData(wh.Conn,cityName)
			fmt.Println(wdat)
			fmt.Println(err)
			if err!=nil{
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			bdata,err:=json.Marshal(wdat)
			if err!=nil{
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}
			w.Write(bdata)
		}else{
			w.WriteHeader(http.StatusMethodNotAllowed)
			m,_:=json.Marshal("method not allowed")
			w.Write(m)
		}
	}else{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println("bad requ")
		m,_:=json.Marshal("no such url")
		w.Write(m)
	}
	
}