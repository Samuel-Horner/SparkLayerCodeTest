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

type DeleteTodoJSON struct {
	Index int `json:"index"`
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

func getHandler(w http.ResponseWriter, _ *http.Request) {
	log.Print("Returning list: ", toDoList)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(toDoList)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// https://pkg.go.dev/encoding/json#example-Decoder
	var decoder *json.Decoder = json.NewDecoder(r.Body)
	var toDo ToDo
	var err error = decoder.Decode(&toDo)

	if err != nil {
		log.Print("Failed to decode request body with error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if toDo.Title == "" || toDo.Description == "" {
		log.Print("Invalid todo item - empty title or description.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	toDoList = append(toDoList, toDo)

	log.Print("Returning item: ", toDo)
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(toDo)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// https://pkg.go.dev/encoding/json#example-Decoder
	var decoder *json.Decoder = json.NewDecoder(r.Body)
	var deleteToDoJSON DeleteTodoJSON
	var err error = decoder.Decode(&deleteToDoJSON)

	if err != nil {
		log.Print("Failed to decode request body with error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if deleteToDoJSON.Index >= len(toDoList) {
		log.Print("Invalid todo item - empty title or description.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Possible issue with memory leaking here
	// See https://go.dev/wiki/SliceTricks discussion on cut/delete
	// Assuming its a non-issue since strings are read only
	toDoList = append(toDoList[:deleteToDoJSON.Index], toDoList[deleteToDoJSON.Index+1:]...)

	log.Print("Deleted todo: ", deleteToDoJSON)
	w.WriteHeader(http.StatusOK)
}

func optionsHandler(w http.ResponseWriter, _ *http.Request) {
	// Needed for the DELETE method to work (need to respond OK with cors settings since DELETE method disabled by default)
	w.WriteHeader(http.StatusOK)
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Method not allowed: ", r.Method)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func ToDoListHandler(w http.ResponseWriter, r *http.Request) {
	// Always respond with these settings
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")

	// Your code here
	log.Print("Recieved ", r.Method, " for ", r.URL.Path, " from ", r.RemoteAddr)

	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	case http.MethodPost:
		postHandler(w, r)
	case http.MethodDelete:
		deleteHandler(w, r)
	case http.MethodOptions:
		optionsHandler(w, r)
	default:
		defaultHandler(w, r)
	}
}
