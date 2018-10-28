package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"golang-poc/app"
	"golang-poc/controllers"
	"net/http"
	"os"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/user", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/user/login", controllers.Authenticate).Methods("POST")

	router.HandleFunc("/auction", controllers.CreateAuction).Methods("POST")
	router.HandleFunc("/auction", controllers.GetAllAuctions).Methods("GET")
	router.HandleFunc("/auction/{id:[0-9]+}", controllers.GetAuctionById).Methods("GET")
	router.HandleFunc("/me/auction", controllers.GetAuctionsForUser).Methods("GET")

	router.HandleFunc("/auction/{id:[0-9]+}/file", controllers.CreateAuctionFile).Methods("POST")
	router.HandleFunc("/auction/{id:[0-9]+}/file/{fileId:[0-9]+}", controllers.GetAuctionFileById).Methods("GET")

	router.Use(app.JwtAuthentication) //attach JWT auth middleware

	router.NotFoundHandler = app.NotFoundHandler

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router) //Launch the app, visit localhost:8000
	if err != nil {
		fmt.Print(err)
	}
}
