package main

import (
	"crypto/sha512"
	_ "expvar"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
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
	random := strconv.Itoa(rand.IntN(1000000))
	hash := h.Sum([]byte(random))
	return string(hash), nil
}

func hashHandler(w http.ResponseWriter, _ *http.Request) {
	hash, err := calcUniqueHash()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hash value: %x", hash)
}

func main() {
	http.HandleFunc("/hash", hashHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
