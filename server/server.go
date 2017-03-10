package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("vim-go")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", Ping).Methods("POST")
	router.HandleFunc("/test", Test).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))

}

// Ping service
func Ping(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	device := r.FormValue("device")
	ip := r.FormValue("ip")
	log.Printf("device=%s, ip=%s", device, ip)
}

// Test service
func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}
