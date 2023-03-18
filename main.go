package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
	"sync"
	"time"
)

var (
	httpHost     string
	httpPort     int
	fingerprints = make(Fingerprints)
	mu           sync.Mutex
)

type Alert struct {
	Fingerprint string    `json:"fingerprint"`
	StartsAt    time.Time `json:"startsAt"`
}

type Data struct {
	Status string  `json:"status"`
	Alerts []Alert `json:"alerts"`
}

// Fingerprints keeps track of the number of alerts received
// by fingerprint and StartsAt time.
type Fingerprints map[string]map[time.Time]int

func updateFingerprints(v Data) {
	mu.Lock()
	defer mu.Unlock()
	for _, alert := range v.Alerts {
		m, ok := fingerprints[alert.Fingerprint]
		if !ok {
			m = make(map[time.Time]int)
		}
		m[alert.StartsAt] += 1
		fingerprints[alert.Fingerprint] = m
	}
}

func parseFlags() {
	flag.StringVar(&httpHost, "http-host", "127.0.0.1", "The HTTP host")
	flag.IntVar(&httpPort, "http-port", 8080, "The HTTP port")
	flag.Parse()
}

func main() {
	parseFlags()
	httpAddr := fmt.Sprintf("%s:%d", httpHost, httpPort)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		v := Data{}
		if err := json.Unmarshal(b, &v); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		updateFingerprints(v)
		log.Println(string(b))
	})
	http.HandleFunc("/fingerprints", func(w http.ResponseWriter, r *http.Request) {
		b, err := func() ([]byte, error) {
			mu.Lock()
			defer mu.Unlock()
			return json.Marshal(fingerprints)
		}()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b)
	})
	log.Printf("Listening on %s\n", httpAddr)
	http.ListenAndServe(httpAddr, nil)
}
