package handler

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/tlbvb/weatherestapi/internal/database"
	"github.com/tlbvb/weatherestapi/openweathermap"
)


type Case struct {
	Method string 
	Path   string
	Query  string
	Status int
	Result interface{}
	Body   interface{}
}

var (
	client = &http.Client{Timeout: time.Second}
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
		('Tokyo', 30.8, 32.5, 28.0, 34.0, 80, 1005);`,

	}

	for _, q := range qs {
		_, err := db.Exec(context.Background(),q)
		if err != nil {
			panic(err)
		}
	}
}



func TestHandler(t *testing.T) {
	conn,err:=database.Connect("postgres://testUser:test@localhost:5432/weather?sslmode=disable")
	if err!=nil{
		t.Fatalf("Couldn't connect to the database",)
	}
	PrepareTestApis(conn)


	ts := httptest.NewServer(&WeatherHandler{Conn: conn})
	fmt.Println(ts)
	cases := []Case{
		Case{
			Path: "/",
			Status: http.StatusBadRequest,
			Result: "no such url",
			
		},
		Case{
			Path: "/weather",
			Status: http.StatusBadRequest,
			Result: "didn't mention the city",
		},
		Case{
			Path:  "/weather",
			Query: "city=Landon",
			Status: http.StatusBadRequest,
			Result: "No such city in table",
		},
		Case{
			Path:  "/weather",
			Query: "city=London",
			Result: openweathermap.WeatherData{Temp:18.2,FeelsLike:  20.5,MinTemp:  15.0,MaxTemp:  22.0, Humidity:65,Pressure: 1018,City: "London"},
		},
		Case{
			Path:  "/weather",
			Query: "city=London",
			Method: http.MethodPut,
		},
		Case{
			Path:  "/weather",
			Query: "city=",
			Method: http.MethodPut,
			Status: http.StatusBadRequest,
			Result: "didn't mention the city",
		},
		Case{
			Path:  "/weather",
			Query: "city=Almaty",
			Method: http.MethodPost,
			Status: http.StatusMethodNotAllowed,
			Result: "method not allowed",
		},
	}
	runCases(t, ts, conn, cases)
}

func runCases(t *testing.T, ts *httptest.Server, db *pgx.Conn, cases []Case) {
	for idx, item := range cases {
		var (
			err      error
			result   interface{}
			expected interface{}
			req      *http.Request
		)
		fmt.Println("testcase",idx,item.Method,item.Query, item.Body )
		caseName := fmt.Sprintf("case %d: [%s] %s %s", idx, item.Method, item.Path, item.Query)

		req, _ = http.NewRequest(item.Method, ts.URL+item.Path+"?"+item.Query, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)

		if item.Status == 0 {
			item.Status = http.StatusOK
		}

		if resp.StatusCode != item.Status {
			t.Fatalf("[%s] expected http status %v, got %v", caseName, item.Status, resp.StatusCode)
			continue
		}
		fmt.Println(body)
		err = json.Unmarshal(body, &result)
		fmt.Println("result",result)
		fmt.Println(body)
		if err != nil {
			t.Fatalf("[%s] cant unpack json: %v", caseName, err)
			continue
		}

		data, _ := json.Marshal(item.Result)
		json.Unmarshal(data, &expected)
		if item.Method==http.MethodPut{
			if result==""{
				t.Fatalf("Incorrect result")
			}
		}else{
			if !reflect.DeepEqual(result, expected) {
				t.Fatalf("[%s] results not match\nGot : %#v\nWant: %#v", caseName, result, expected)
				continue
			}
		}
	}

}
