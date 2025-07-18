package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Adding a struct to represent the hardware device
type Device struct {
	ID       int    `json:"id"`
	Mac      string `json:"mac"`
	Firmware string `json:"firmware"`
}

// Declare a struct for the first Gauge metric with the number of connected devices
type metrics struct {
	// - prometheus.Gauge represents a single numerical value for devices. Used for the same devices.
	// - prometheus.GaugeVector is a bundle of Gauges with the same name but different label. Use when
	//     when you have several devices types.
	devices prometheus.Gauge
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

// Metric function
// Input: prometheus.Registerer, Output: *metrics
func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",             // Namespace is a unique word prefix
			Name:      "connected_devices", // Metric name following Prometheus conventions
			Help:      "Number of currently connected devices.",
		}),
	}
	reg.MustRegister(m.devices) // Registering with Prometheus registry
	return m                    // Return pointer
}

// Add a handler to send the connected devices through the endpoint
// The star is a pointer
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

func main() {
	reg := prometheus.NewRegistry() // Create a non-global registry without any collectors
	m := NewMetrics(reg)            // Create metrics using the NewMetrics function

	m.devices.Set(float64(len(dvs))) // Set de Gauge according to the connected devices

	// Custom Prometheus handler with the new registry
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	// HTTP endpoint using the custom handler
	http.Handle("/metrics", promHandler)
	// Enable the /devices endpoint passing the function itself,
	// not a call to it using HandleFunc
	http.HandleFunc("/devices", getDevices)
	http.ListenAndServe(":2112", nil)
}
