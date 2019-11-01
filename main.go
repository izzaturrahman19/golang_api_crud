package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

var ctx= context.Background()

type student struct{

	Name string `bson:"name"`
	Grade int   `bson:"Grade"`
	Category string   `bson:"Category"`

}

func connect()(*mongo.Database, error){
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil,err
	}

	return client.Database("belajar_golang"), nil
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	if r.Method== "GET"{

		var tmpl = template.Must(template.New("form").ParseFiles("view.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	
	if r.Method== "POST"{

	db, err :=connect()
	if err != nil {
        log.Fatal(err.Error())
	}

	var tmpl = template.Must(template.New("result").ParseFiles("view.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	var name = r.FormValue("name")
	var bubble = r.Form.Get("grade")
	var grade int
	if grade, err := strconv.ParseInt(bubble, 10, 32); err == nil {
		fmt.Printf("%T, %v\n", grade, grade)
	}
	var category = r.Form.Get("category")
	
	
	var data = map[string]string{"name": name, "category": category }
	var angka = map[string]int{ "grade": grade}
	
	_, err = db.Collection("student").InsertOne(ctx, student{name, grade, category})
	if err != nil {
		log.Fatal(err.Error())
	}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		if err := tmpl.Execute(w, angka); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return

}}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	
	if r.Method== "PATCH"{

	db, err := connect()
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = db.Collection("student").UpdateOne(ctx, selector, bson.M{"$set": changes})
	if err != nil {
		log.Fatal(err.Error())
	}

}}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/process", createEvent).Methods("POST")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/update}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":6060", router))
}