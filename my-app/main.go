package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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
	info    *prometheus.GaugeVec // Metric for metadata
}

// Global vars
var dvs []Device   // Hold all the connected devices
var version string // App version

// Use the init function to define some values required for demo
func init() {
	// Hardcoded version for demo
	version = "2.10.5"

	// Devices values during application starting
	dvs = []Device{
		{1, "5F-22-CC-1F-43-82", "2.1.6"},
		{2, "EF-2B-C4-F5-D6-34", "2.1.6"},
	}
}

// Metric function
// Input: prometheus.Registerer, Output: *metrics
func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		// First metric configured for amount of devices connected
		devices: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "myapp",                                  // Namespace is a unique word prefix
			Name:      "connected_devices",                      // Metric name following Prometheus conventions
			Help:      "Number of currently connected devices.", // Metric description
		}),
		// Second metric configured for version number
		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "myapp",                                     // Namespace is a unique word prefix
			Name:      "info",                                      // Metric name following Prometheus conventions
			Help:      "Information about the My App environment.", // Metric description
		},
			[]string{"version"}),
	}
	reg.MustRegister(m.devices, m.info) // Registering every metric with Prometheus registry
	return m                            // Return pointer
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
	reg := prometheus.NewRegistry()               // Create a non-global registry without any collectors
	reg.MustRegister(collectors.NewGoCollector()) // Add the included collector to the registry
	m := NewMetrics(reg)                          // Create metrics using the NewMetrics function

	m.devices.Set(float64(len(dvs)))                          // Set de Gauge according to the connected devices
	m.info.With(prometheus.Labels{"version": version}).Set(1) // If version is empty, uses 1 as default

	// Custom Prometheus handler with the new registry
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	// HTTP endpoint using the custom handler
	http.Handle("/metrics", promHandler)
	// Enable the /devices endpoint passing the function itself,
	// not a call to it using HandleFunc
	http.HandleFunc("/devices", getDevices)
	http.ListenAndServe(":2112", nil)
}
