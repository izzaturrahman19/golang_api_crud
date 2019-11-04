package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

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

func readEvent(w http.ResponseWriter, r *http.Request) {
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
	if r.Method== "GET"{
		var tmpl = template.Must(template.New("form").ParseFiles("view.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
	}
	
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
	grade, err := strconv.Atoi(bubble)
	if err == nil {
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

	}
	http.Error(w, "", http.StatusBadRequest)
}


func deleteEvent(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET"{
		
		var tmpl = template.Must(template.New("delete").ParseFiles("view.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return
	}
	if r.Method == "POST"{

		db, err := connect()
    if err != nil {
        log.Fatal(err.Error())
	}
	
    var tmpl = template.Must(template.New("result_d").ParseFiles("view.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
	var vector = r.FormValue("data_d")
	var selector = bson.M{"name": vector}

	var data = map[string]string{"data": vector}

	_, err = db.Collection("student").DeleteOne(ctx, selector)
	if err != nil {
        log.Fatal(err.Error())
	}
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println("Remove Success !")
return
	}
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method== "GET"{
		var tmpl = template.Must(template.New("update").ParseFiles("view.html"))
		var err = tmpl.Execute(w, nil)

		if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
        return	
	}

	if r.Method== "POST"{

		db, err :=connect()
	if err != nil {
        log.Fatal(err.Error())
	}

	var tmpl = template.Must(template.New("result_u").ParseFiles("view.html"))

		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	var vector = r.FormValue("data_u")
	var selector = bson.M{"name": vector}

	var name = r.FormValue("name_u")
	var bubble = r.Form.Get("grade_u")
	grade, err := strconv.Atoi(bubble)
	if err == nil {
		fmt.Printf("%T, %v\n", grade, grade)
	}
	var category = r.Form.Get("category_u")
	var changes = student{name, grade, category}

	var data = map[string]string{"name": name, "category": category }
	var angka = map[string]int{ "grade": grade}
	
	_, err = db.Collection("student").UpdateOne(ctx, selector, bson.M{"$set": changes})
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
	}
	http.Error(w, "", http.StatusBadRequest)
}



func main(){

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/read", readEvent)
	router.HandleFunc("/event", createEvent).Methods("GET")
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/update", updateEvent).Methods("GET")
	router.HandleFunc("/update", updateEvent).Methods("POST")
	router.HandleFunc("/delete", deleteEvent).Methods("GET")
	router.HandleFunc("/delete", deleteEvent).Methods("POST")
	log.Fatal(http.ListenAndServe(":6060", router))
}
