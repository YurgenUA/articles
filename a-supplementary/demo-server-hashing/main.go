package main

import (
	"crypto/sha512"
	_ "expvar"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"runtime"
	"strconv"
)

func loadFile() ([]byte, error) {
	body, err := os.ReadFile("./assets/minikonda3.pkg")
	if err != nil {
		return nil, err
	}
	return body, nil
}

func calcUniqueHash() (string, error) {
	body, err := loadFile()
	if err != nil {
		return "", err
	}
	h := sha512.New()
	h.Write(body)
	random := strconv.Itoa(rand.Intn(1000000))
	hash := h.Sum([]byte(random))
	return string(hash), nil
}

func memstats2() any {
	log.Println("Collecting memory stats")
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	log.Printf("Memory stats collected: %+v", *stats)
	return *stats
}

func hashHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    log.Printf("Started %s %s", r.Method, r.URL.Path)

    hash, err := calcUniqueHash()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Printf("Completed %s %s with error in %v", r.Method, r.URL.Path, time.Since(start))
        return
    }
    fmt.Fprintf(w, "Hash value: %x", hash)
    log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
}

func main() {
	http.HandleFunc("/hash", hashHandler)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
