package main

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// Adding a struct to represent the hardware device
type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

// Declare a global var to hold all the connected devices
var dvs []Device

// Use the init function to define some devices values during application starting
func init() {
	dvs = []Device{
		{1, "5F-22-CC-1F-43-82", "2.1.6"},
		{2, "EF-2B-C4-F5-D6-34", "2.1.6"},
	}
}

func main() {
	// HTTP endpoint using the promhttp library
	http.Handle("/metrics", promhttp.Handler())
	// Enable the /devices endpoint passing the function itself,
	// not a call to it using HandleFunc
	http.HandleFunc("/devices", getDevices)
	http.ListenAndServe(":2112", nil)
}

// Add a handler to send the connected devices through the endpoint
func getDevices(w http.ResponseWriter, r *http.Request) {
	// json.Marshal converts go structs into json strings
	b, err := json.Marshal(dvs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Define the content header as "application/json"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Write the connection
	w.Write(b)
}
