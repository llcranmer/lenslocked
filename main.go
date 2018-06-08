package main

import (
	"./controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	usersC := controllers.NewUsers()
	staticC := controllers.NewStatic()

	r := mux.NewRouter()

	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/contact", staticC.ContactView).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}
