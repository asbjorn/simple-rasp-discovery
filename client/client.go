package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"

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
	Port   string `short:"p" long:"port" description:"Port number if needed. Default is 80" default:"80"`
}

var options Options
var parser = flags.NewParser(&options, flags.Default)

// GetPayload creates the data
func GetPayload() url.Values {
	myip := GetMyIP()

	data := url.Values{}
	data.Add("device", "happymeter")
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

	discoveryService := options.Server

	u, _ := url.ParseRequestURI(discoveryService)
	urlStr := fmt.Sprintf("%v", u)

	log.Println("Trying to connect to", urlStr)

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
