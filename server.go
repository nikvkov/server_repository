package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Sity struct {
	Id      int
	City    string
	Country string
}

var database *sql.DB

type msg string

func (m msg) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprint(resp, m)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	cnt := vars["cnt"]

	var query1 string = fmt.Sprintf("select cities.id, city, country from  cities join regions on region_id = regions.id join countries on regions.country_id = countries.id where country LIKE '%s' order by city", cnt)
	//var query1 string = "select cities.id, city, country from  cities join regions on region_id = regions.id join countries on regions.country_id = countries.id where country LIKE 'Украина' order by city"

	//fmt.Fprint(w, query1)
	rows, err := database.Query(query1)

	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	sities := []Sity{}

	for rows.Next() {
		p := Sity{}
		err := rows.Scan(&p.Id, &p.City, &p.Country)
		if err != nil {
			fmt.Println(err)
			continue
		}
		sities = append(sities, p)
	}

	//tmpl, _ := template.ParseFiles("templates/index.html")
	//tmpl.Execute(w, sities)

	sliceVar2, _ := json.Marshal(sities)

	fmt.Fprint(w, string(sliceVar2))
	//msgHandler := msg("string(sliceVar2)")
	//	http.ListenAndServe("localhost:8181", msgHandler)
	//fmt.Println(string(sliceVar2))
	//http.ListenAndServe("localhost:8181", msgHandler)
}

func main() {

	db, err := sql.Open("mysql", "root:developer@/skillful_hands")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()
	//http.HandleFunc("/country/{cnt}", IndexHandler)
	router := mux.NewRouter()
	router.HandleFunc("/country/{cnt}", IndexHandler)
	http.Handle("/", router)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)

}
