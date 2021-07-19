package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

//TICKETS
type Ticket struct {
	Id       int    `json:"id"`
	Issue    string `json:"issue"`
	Priority string `json:"priority"`
	Status   string `json:"status"`
}

//DB CONNECTION
func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "admin:password@/ticketsDB")
	if err != nil {
		panic(err.Error())
	}
	return db
}

//CREATE TICKET
func createTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		reqBody, _ := ioutil.ReadAll(r.Body)
		var ticket Ticket
		json.Unmarshal(reqBody, &ticket)

		db := dbConn()
		issue := ticket.Issue
		priority := ticket.Priority
		status := ticket.Status

		insTicket, err := db.Prepare("INSERT INTO tickets(issue, priority, status) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}

		insTicket.Exec(issue, priority, status)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(ticket); err != nil {
			panic(err)
		}

		defer db.Close()
	}
}

//RETURN ALL TICKETS
func returnAllTickets(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	selTickets, err := db.Query("SELECT * FROM tickets")
	if err != nil {
		panic(err.Error())
	}

	ticket := Ticket{}
	allTickets := []Ticket{}
	for selTickets.Next() {
		var id int
		var issue, priority, status string
		err = selTickets.Scan(&id, &issue, &priority, &status)
		if err != nil {
			panic(err.Error())
		}
		ticket.Id = id
		ticket.Issue = issue
		ticket.Priority = priority
		ticket.Status = status
		allTickets = append(allTickets, ticket)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(allTickets); err != nil {
		panic(err)
	}

	defer db.Close()
}

//RETURN TICKET
func returnTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		vars := mux.Vars(r)
		id := vars["id"]

		db := dbConn()
		selTicket, err := db.Query("SELECT * FROM tickets WHERE id=?", id)
		if err != nil {
			panic(err.Error())
		}

		ticket := Ticket{}
		for selTicket.Next() {
			var id int
			var issue, priority, status string
			err = selTicket.Scan(&id, &issue, &priority, &status)
			if err != nil {
				panic(err.Error())
			}
			ticket.Id = id
			ticket.Issue = issue
			ticket.Priority = priority
			ticket.Status = status
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(ticket); err != nil {
			panic(err)
		}

		defer db.Close()
	}
}

//UPDATE TICKET
func updateTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method == "PUT" {

		reqBody, _ := ioutil.ReadAll(r.Body)
		var ticket Ticket
		json.Unmarshal(reqBody, &ticket)

		db := dbConn()

		id := ticket.Id
		issue := ticket.Issue
		priority := ticket.Priority
		status := ticket.Status

		updateTicket, err := db.Prepare("UPDATE tickets SET issue=?, priority=?, status=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		updateTicket.Exec(issue, priority, status, id)

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(ticket); err != nil {
			panic(err)
		}

		defer db.Close()
	}
}

//DELETE TICKET
func deleteTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method == "DELETE" {

		vars := mux.Vars(r)
		id := vars["id"]

		db := dbConn()
		delTicket, err := db.Prepare("DELETE FROM tickets WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		delTicket.Exec(id)
		defer db.Close()
	}
}

func main() {

	log.Println("Server started on: localhost")
	router := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"X-Request-Widht", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/", returnAllTickets).Methods("GET")
	router.HandleFunc("/ticket", createTicket).Methods("POST")
	router.HandleFunc("/ticket/{id}", returnTicket).Methods("GET")
	router.HandleFunc("/ticket/update", updateTicket).Methods("PUT")
	router.HandleFunc("/ticket/{id}", deleteTicket).Methods("DELETE")

	err := http.ListenAndServe(":8080", handlers.CORS(headers, methods, origins)(router))
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
