package main

import (
	"fmt"
	"net/http"

	"github.com/tlbvb/weatherestapi/internal/database"
	"github.com/tlbvb/weatherestapi/internal/handler"
)



func main(){
	host:=":8080"
	connStr := "postgres://testUser:test@localhost:5432/weather?sslmode=disable"
	conn,err:=database.Connect(connStr)
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println("connected")	
	http.ListenAndServe(host,&handler.WeatherHandler{Conn:conn})
}


