package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/shirou/gopsutil/host"
)

// GetMyIP returns the local IP
func GetMyIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
		return ""
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			return ipnet.IP.String()
		}
	}
	return "n/a"
}

// Options flags
type Options struct {
	Server string `short:"s" long:"server" description:"Full http address for the server to connect to. i.e. http://example.com:8080/my_endpoint" required:"true"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

// GetPayload creates the data
func GetPayload() url.Values {
	myip := GetMyIP()

	data := url.Values{}
	hostname, err := os.Hostname()
	if err != nil {
		data.Add("device", "unknown")
	} else {
		data.Add("device", hostname)
	}

	data.Add("ip", myip)

	uptime, err := host.Uptime()
	log.Println("System uptime:", uptime)
	if err != nil {
		log.Fatalln(err)
	}
	uptimeStr := strconv.Itoa(int(uptime))
	data.Add("uptime", uptimeStr)

	return data
}

func main() {
	if _, err := parser.Parse(); err != nil {
		return
	}

	current_time := time.Now().Local()
	log.Println("Raspberry discovery client: ", current_time)

	discoveryService := options.Server

	u, _ := url.ParseRequestURI(discoveryService)
	urlStr := fmt.Sprintf("%v", u)

	log.Println("Connecting to", urlStr)

	data := GetPayload()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Println(resp.Status)
}
