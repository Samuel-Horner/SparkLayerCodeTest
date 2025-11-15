package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ToDo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ToDoList []ToDo

const port int = 8080

var toDoList ToDoList = ToDoList{}

func main() {
	// Your code here

	// https://pkg.go.dev/net/http#example-HandleFunc
	http.HandleFunc("/", ToDoListHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Returning list: ", toDoList)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(toDoList)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var decoder *json.Decoder = json.NewDecoder(r.Body)
	var toDo ToDo
	var err error = decoder.Decode(&toDo)

	if err != nil {
		log.Print("Failed to decode request body with error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	toDoList = append(toDoList, toDo)

	log.Print("Returning item: ", toDo)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(toDo)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Method not allowed: ", r.Method)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func ToDoListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Your code here
	log.Print("Recieved ", r.Method, " for ", r.URL.Path, " from ", r.RemoteAddr)

	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)
	default:
		defaultHandler(w, r)
	}
}
