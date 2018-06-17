package main

import (
	"./controllers"
	"./models"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "kbabkin"
	dbname = "postgres"
)

func main() {
	PSQLInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)

	us, err := models.NewUserService(PSQLInfo)
	if err != nil {
		panic(err)
	}

	defer us.Close()

	usersC := controllers.NewUsers(us)
	staticC := controllers.NewStatic()

	//us.DestructiveReset()
	us.AutoMigrate()

	r := mux.NewRouter()

	r.Handle("/", staticC.HomeView).Methods("GET")
	r.Handle("/contact", staticC.ContactView).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	http.ListenAndServe(":3000", r)
}
