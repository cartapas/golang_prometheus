package main

import (
	"encoding/json"
	"log"
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

type registerDevicesHandler struct {
	metrics *metrics // Property to increment the device creation
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
		{3, "35-2B-C4-F5-17-84", "2.1.7"},
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

// Register a new device using a new handler because the standard one didn't manage a switch between
// GET, POST methods
// Rename to ServeHTTP and add the registry device handler (rdh)
func (rdh registerDevicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getDevices(w, r)
	case "POST":
		createDevice(w, r, rdh.metrics) // Adding arguments to the createDevice func
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Add a new device to the current list
func createDevice(w http.ResponseWriter, r *http.Request, m *metrics) {
	var dv Device

	err := json.NewDecoder(r.Body).Decode(&dv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dvs = append(dvs, dv)

	// The Inc could be used to increment the counter, but other option will be used to be more accurate
	// Put `after dvs = append(dvs, dv) to update the counter
	// m.devices.Inc()
	m.devices.Set(float64(len(dvs)))

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Device created!"))
}

func main() {
	reg := prometheus.NewRegistry()               // Create a non-global registry without any collectors
	reg.MustRegister(collectors.NewGoCollector()) // Add the included collector to the registry
	m := NewMetrics(reg)                          // Create metrics using the NewMetrics function

	m.devices.Set(float64(len(dvs)))                          // Set de Gauge according to the connected devices
	m.info.With(prometheus.Labels{"version": version}).Set(1) // If version is empty, uses 1 as default

	// Use multiplexers to run several servers avoiding the external endpoints exposicion
	// dMux is used for the customized endpoints
	dMux := http.NewServeMux()
	// Initialize the rdh with a pointer to the metrics
	rdh := registerDevicesHandler{metrics: m}
	// Replace the basic http.HandleFunc with the switch and add the created rdh
	dMux.Handle("/devices", rdh)

	// pMux is used for the default `/metrics endpoint
	pMux := http.NewServeMux()
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	pMux.Handle("/metrics", promHandler)

	// Listen for customized `devices`endpoint
	go func() {
		log.Fatal(http.ListenAndServe(":8080", dMux))
	}()

	// Listen for default `/metrics`` endpoint
	go func() {
		log.Fatal(http.ListenAndServe(":8081", pMux))
	}()

	// Keep the app running
	select {}
}
