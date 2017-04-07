package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var port = 8080

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", Ping).Methods("POST")
	router.HandleFunc("/test", Test).Methods("GET")

	log.Println("Raspberry Discovery Server started on port..: " + strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))

}

// Ping service
func Ping(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	device := r.FormValue("device")
	ip := r.FormValue("ip")
	log.Printf("device=%s, ip=%s", device, ip)

	for k, v := range r.Form {
		log.Printf("debug: key=%s, value=%s", k, v)
	}
}

// Test service
func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}
