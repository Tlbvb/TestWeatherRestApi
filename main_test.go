package main

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/tlbvb/weatherestapi/internal/database"
	"github.com/tlbvb/weatherestapi/internal/handler"
)

// CaseResponse
type CR map[string]interface{}

type Case struct {
	Method string // GET по-умолчанию в http.NewRequest если передали пустую строку
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

func CleanupTestApis(db *pgx.Conn) {
	qs := []string{
		`DROP TABLE IF EXISTS items;`,
		`DROP TABLE IF EXISTS users;`,
	}
	for _, q := range qs {
		_, err := db.Exec(context.Background(),q)
		if err != nil {
			panic(err)
		}
	}
}

func TestApis(t *testing.T) {
	//conn,err:=pgx.Connect(context.Background(),"postgres://testUser:test@localhost:5432/weather?sslmode=disable")
	// conn,err:=database.Connect("postgres://testUser:test@localhost:5432/weatherr?sslmode=disable")
	// if err!=nil{
	// 	t.Fatalf("Couldn't connect to the database",)
	// }
	conn,err:=database.Connect("postgres://testUser:test@localhost:5432/weather?sslmode=disable")
	if err!=nil{
		t.Fatalf("Couldn't connect to the database",)
	}
	PrepareTestApis(conn)

	defer CleanupTestApis(conn)


	ts := httptest.NewServer(&handler.WeatherHandler{Conn: conn})
	fmt.Println(ts)
	cases := []Case{
		Case{
			Path: "/", // список таблиц
			Status: http.StatusBadRequest,
			Result: "no such url",
			// Result: CR{
			// 	"response": CR{
			// 		"tables": []string{"items", "users"},
			// 	},
			// },
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
	// 	Case{
	// 		Path:  "/items",
	// 		Query: "limit=1&offset=1",
	// 		Result: CR{
	// 			"response": CR{
	// 				"records": []CR{
	// 					CR{
	// 						"id":          2,
	// 						"title":       "memcache",
	// 						"description": "Рассказать про мемкеш с примером использования",
	// 						"updated":     nil,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path: "/items/1",
	// 		Result: CR{
	// 			"response": CR{
	// 				"record": CR{
	// 					"id":          1,
	// 					"title":       "database/sql",
	// 					"description": "Рассказать про базы данных",
	// 					"updated":     "rvasily",
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path:   "/items/100500",
	// 		Status: http.StatusNotFound,
	// 		Result: CR{
	// 			"error": "record not found",
	// 		},
	// 	},

	// 	// тут идёт создание и редактирование
	// 	Case{
	// 		Path:   "/items/",
	// 		Method: http.MethodPut,
	// 		Body: CR{
	// 			"id":          42, // auto increment primary key игнорируется при вставке
	// 			"title":       "db_crud",
	// 			"description": "",
	// 		},
	// 		Result: CR{
	// 			"response": CR{
	// 				"id": 3,
	// 			},
	// 		},
	// 	},
	// 	// это пример хрупкого теста
	// 	// если много раз вызывать один и тот же тест - записи будут добавляться
	// 	// поэтому придётся сделать сброс базы каждый раз в PrepareTestData
	// 	Case{
	// 		Path: "/items/3",
	// 		Result: CR{
	// 			"response": CR{
	// 				"record": CR{
	// 					"id":          3,
	// 					"title":       "db_crud",
	// 					"description": "",
	// 					"updated":     nil,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path:   "/items/3",
	// 		Method: http.MethodPost,
	// 		Body: CR{
	// 			"description": "Написать программу db_crud",
	// 		},
	// 		Result: CR{
	// 			"response": CR{
	// 				"updated": 1,
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path: "/items/3",
	// 		Result: CR{
	// 			"response": CR{
	// 				"record": CR{
	// 					"id":          3,
	// 					"title":       "db_crud",
	// 					"description": "Написать программу db_crud",
	// 					"updated":     nil,
	// 				},
	// 			},
	// 		},
	// 	},

	// 	// обновление null-поля в таблице
	// 	Case{
	// 		Path:   "/items/3",
	// 		Method: http.MethodPost,
	// 		Body: CR{
	// 			"updated": "autotests",
	// 		},
	// 		Result: CR{
	// 			"response": CR{
	// 				"updated": 1,
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path: "/items/3",
	// 		Result: CR{
	// 			"response": CR{
	// 				"record": CR{
	// 					"id":          3,
	// 					"title":       "db_crud",
	// 					"description": "Написать программу db_crud",
	// 					"updated":     "autotests",
	// 				},
	// 			},
	// 		},
	// 	},

	// 	// обновление null-поля в таблице
	// 	Case{
	// 		Path:   "/items/3",
	// 		Method: http.MethodPost,
	// 		Body: CR{
	// 			"updated": nil,
	// 		},
	// 		Result: CR{
	// 			"response": CR{
	// 				"updated": 1,
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path: "/items/3",
	// 		Result: CR{
	// 			"response": CR{
	// 				"record": CR{
	// 					"id":          3,
	// 					"title":       "db_crud",
	// 					"description": "Написать программу db_crud",
	// 					"updated":     nil,
	// 				},
	// 			},
	// 		},
	// 	},

	// 	// ошибки
	// 	Case{
	// 		Path:   "/items/3",
	// 		Method: http.MethodPost,
	// 		Status: http.StatusBadRequest,
	// 		Body: CR{
	// 			"id": 4, // primary key нельзя обновлять у существующей записи
	// 		},
	// 		Result: CR{
	// 			"error": "field id have invalid type",
	// 		},
	// 	},
	// 	Case{
	// 		Path:   "/items/3",
	// 		Method: http.MethodPost,
	// 		Status: http.StatusBadRequest,
	// 		Body: CR{
	// 			"title": 42,
	// 		},
	// 		Result: CR{
	// 			"error": "field title have invalid type",
	// 		},
	// 	},
	// 	Case{
	// 		Path:   "/items/3",
	// 		Method: http.MethodPost,
	// 		Status: http.StatusBadRequest,
	// 		Body: CR{
	// 			"title": nil,
	// 		},
	// 		Result: CR{
	// 			"error": "field title have invalid type",
	// 		},
	// 	},

	// 	Case{
	// 		Path:   "/items/3",
	// 		Method: http.MethodPost,
	// 		Status: http.StatusBadRequest,
	// 		Body: CR{
	// 			"updated": 42,
	// 		},
	// 		Result: CR{
	// 			"error": "field updated have invalid type",
	// 		},
	// 	},

	// 	// удаление
	// 	Case{
	// 		Path:   "/items/3",
	// 		Method: http.MethodDelete,
	// 		Result: CR{
	// 			"response": CR{
	// 				"deleted": 1,
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path:   "/items/3",
	// 		Method: http.MethodDelete,
	// 		Result: CR{
	// 			"response": CR{
	// 				"deleted": 0,
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path:   "/items/3",
	// 		Status: http.StatusNotFound,
	// 		Result: CR{
	// 			"error": "record not found",
	// 		},
	// 	},

	// 	// и немного по другой таблице
	// 	Case{
	// 		Path: "/users/1",
	// 		Result: CR{
	// 			"response": CR{
	// 				"record": CR{
	// 					"user_id":  1,
	// 					"login":    "rvasily",
	// 					"password": "love",
	// 					"email":    "rvasily@example.com",
	// 					"info":     "none",
	// 					"updated":  nil,
	// 				},
	// 			},
	// 		},
	// 	},

	// 	Case{
	// 		Path:   "/users/1",
	// 		Method: http.MethodPost,
	// 		Body: CR{
	// 			"info":    "try update",
	// 			"updated": "now",
	// 		},
	// 		Result: CR{
	// 			"response": CR{
	// 				"updated": 1,
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path: "/users/1",
	// 		Result: CR{
	// 			"response": CR{
	// 				"record": CR{
	// 					"user_id":  1,
	// 					"login":    "rvasily",
	// 					"password": "love",
	// 					"email":    "rvasily@example.com",
	// 					"info":     "try update",
	// 					"updated":  "now",
	// 				},
	// 			},
	// 		},
	// 	},
	// 	// ошибки
	// 	Case{
	// 		Path:   "/users/1",
	// 		Method: http.MethodPost,
	// 		Status: http.StatusBadRequest,
	// 		Body: CR{
	// 			"user_id": 1, // primary key нельзя обновлять у существующей записи
	// 		},
	// 		Result: CR{
	// 			"error": "field user_id have invalid type",
	// 		},
	// 	},
	// 	//не забываем про sql-инъекции
	// 	Case{
	// 		Path:   "/users/",
	// 		Method: http.MethodPut,
	// 		Body: CR{
	// 			"user_id":    2,
	// 			"login":      "qwerty'",
	// 			"password":   "love\"",
	// 			"unkn_field": "love",
	// 		},
	// 		Result: CR{
	// 			"response": CR{
	// 				"user_id": 2,
	// 			},
	// 		},
	// 	},
	// 	Case{
	// 		Path: "/users/2",
	// 		Result: CR{
	// 			"response": CR{
	// 				"record": CR{
	// 					"user_id":  2,
	// 					"login":    "qwerty'",
	// 					"password": "love\"",
	// 					"email":    "",
	// 					"info":     "",
	// 					"updated":  nil,
	// 				},
	// 			},
	// 		},
	// 	},
	// 	// тут тоже возможна sql-инъекция
	// 	// если пришло не число на вход - берём дефолтное значене для лимита-оффсета
	// 	Case{
	// 		Path:  "/users",
	// 		Query: "limit=1'&offset=1\"",
	// 		Result: CR{
	// 			"response": CR{
	// 				"records": []CR{
	// 					CR{
	// 						"user_id":  1,
	// 						"login":    "rvasily",
	// 						"password": "love",
	// 						"email":    "rvasily@example.com",
	// 						"info":     "try update",
	// 						"updated":  "now",
	// 					},
	// 					CR{
	// 						"user_id":  2,
	// 						"login":    "qwerty'",
	// 						"password": "love\"",
	// 						"email":    "",
	// 						"info":     "",
	// 						"updated":  nil,
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
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

		// если у вас случилась это ошибка - значит вы не делаете где-то rows.Close и у вас текут соединения с базой
		// если такое случилось на первом тесте - значит вы не закрываете коннект где-то при иницаилизации в NewDbExplorer
		// if db.Stats().OpenConnections != 1 {
		// 	t.Fatalf("[%s] you have %d open connections, must be 1", caseName, db.Stats().OpenConnections)
		// }

		if item.Method == "" || item.Method == http.MethodGet {
			req, err = http.NewRequest(item.Method, ts.URL+item.Path+"?"+item.Query, nil)
		} else {
			data, err := json.Marshal(item.Body)
			if err != nil {
				panic(err)
			}
			reqBody := bytes.NewReader(data)
			req, err = http.NewRequest(item.Method, ts.URL+item.Path, reqBody)
			req.Header.Add("Content-Type", "application/json")
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		// fmt.Printf("[%s] body: %s\n", caseName, string(body))
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

		// reflect.DeepEqual не работает если нам приходят разные типы
		// а там приходят разные типы (string VS interface{}) по сравнению с тем что в ожидаемом результате
		// этот маленький грязный хак конвертит данные сначала в json, а потом обратно в interface - получаем совместимые результаты
		// не используйте это в продакшен-коде - надо явно писать что ожидается интерфейс или использовать другой подход с точным форматом ответа
		data, err := json.Marshal(item.Result)
		json.Unmarshal(data, &expected)

		if !reflect.DeepEqual(result, expected) {
			t.Fatalf("[%s] results not match\nGot : %#v\nWant: %#v", caseName, result, expected)
			continue
		}
	}

}
