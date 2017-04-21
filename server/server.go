package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jessevdk/go-flags"
)

// Options flags
type Options struct {
	Port int `short:"p" long:"port" description:"Which port number to run this server on. Recommended is 8080." required:"true"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

// Device is a data holder for Device info
type Device struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Uptime string `json:"uptime"`
	Time   string `json:"time"`
}

// Response is the JSON response with all registered Devices included
type Response struct {
	Devices []Device `json:"devices"`
}

var devices = make(map[string]Device)

func main() {
	if _, err := parser.Parse(); err != nil {
		return
	}

	var port = options.Port

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", Ping).Methods("POST")
	router.HandleFunc("/devices", Devices).Methods("GET")

	log.Println("Raspberry Discovery Server started on port..: " + strconv.Itoa(port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))

}

// Ping service
func Ping(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	dName := r.FormValue("device")
	ip := r.FormValue("ip")
	uptime := r.FormValue("uptime")

	if _, ok := devices[dName]; ok {
		log.Println("Updating existing info")
	} else {
		log.Println("New device! ")
	}

	d := devices[dName]
	d.Name = dName
	d.IP = ip
	d.Uptime = uptime
	t := time.Now()
	d.Time = t.Format("2006-01-02 15:04:05")

	devices[dName] = d

	for k, v := range r.Form {
		log.Printf("debug: key=%s, value=%s", k, v)
	}
}

// Devices service
func Devices(w http.ResponseWriter, r *http.Request) {
	log.Println("List devices")
	var units []Device

	for _, val := range devices {
		units = append(units, val)
	}

	respObj := &Response{Devices: units}
	resp, err := json.Marshal(respObj)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		//return 500, err
		log.Fatalf("Facking error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	//fmt.Fprintf(w, string(resp))
	w.Write(resp)
}
