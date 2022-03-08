package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/sairan-ds/go-nats-steaming-project/internal/config"
	"github.com/sairan-ds/go-nats-steaming-project/internal/database"
	"github.com/sairan-ds/go-nats-steaming-project/internal/streaming"
	"github.com/sairan-ds/go-nats-steaming-project/views"
)

var (
	homeView *views.View
	csh map[string]database.Order
)



func main() {
	config.ConfigSetup()
	db := database.SetUp()
	csh = db.Cache
	log.Println("Setting up config Cache")
	s := streaming.Subscribe(db)



	homeView = views.NewView("bootstrap","views/home.gohtml")

	r := mux.NewRouter()
	r.HandleFunc("/{id}", home)
	go func() {
		http.ListenAndServe(":3000", r)
	}()

	signalChan := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			s.Unsubscribe()
			db.Db.Close()
			done <- true
		}
	}()
	<-done
}



func home(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	w.Header().Set("Content-Type", "text/html")
	data := csh[id]
	err := homeView.Template.ExecuteTemplate(w, homeView.Layout, data)
	if err != nil {
		panic(err)
	}
}