package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	water		int
	waterStatus	string
	wind		int
	windStatus	string
	mutex		sync.Mutex
)

type StatusData map[string]int
type KeteranganStatus map[string]string

type Data struct {
	Status StatusData `json:"status"`
	Keterangan KeteranganStatus `json:"keterangan"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	go generateRandomNumber()

	http.HandleFunc("/", handleEvents)

	log.Println("Server started on :5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	status := StatusData{
		"water": water,
		"wind": wind,
	}
	keterangan := KeteranganStatus{
		"water": waterStatus,
		"wind": windStatus,
	}
	data := Data{Status: status, Keterangan: keterangan}
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func generateRandomNumber() {
	for {
		mutex.Lock()
		water = rand.Intn(101)
		if water <= 5 {
			waterStatus = "aman"
		}
		if water >= 6 && water <= 8 {
			waterStatus = "siaga"
		}
		if water > 8 {
			waterStatus = "bahaya"
		}

		wind = rand.Intn(101)
		if wind <= 6 {
			windStatus = "aman"
		}
		if wind >= 7 && wind <= 15 {
			windStatus = "siaga"
		}
		if wind > 15 {
			windStatus = "bahaya"
		}
		mutex.Unlock()

		time.Sleep(15 * time.Second)
	}
}