package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Status struct {
	Status      string `json:"status"`
	Advertizing string `json:"advertizing"`
}

type RangeDetails struct {
	Make  string `json:"make"`
	Model string `json:"model"`
	Link  string `json:"link"`
}

var allRanges = []RangeDetails{
	{
		Make:  "LG",
		Model: "LSGL6337F",
		Link:  "https://www.homedepot.com/p/315878151",
	},
	{
		Make:  "Maytag",
		Model: "MGR6600FZ",
		Link:  "https://www.homedepot.com/p/301061491",
	},
	{
		Make:  "Whirlpool",
		Model: "WFG505M0BS",
		Link:  "https://www.homedepot.com/p/205079331",
	},
	{
		Make:  "Frigidaire",
		Model: "FGGH3047VF",
		Link:  "https://www.homedepot.com/p/309565237",
	},
	{
		Make:  "Samsung",
		Model: "HR1124G",
		Link:  "https://www.homedepot.com/p/315493149",
	},
}

func main() {
	for i, r := range allRanges {
		go func(port int, selectedRange *RangeDetails) {
			http.ListenAndServe(fmt.Sprintf("localhost:808%d", port), setupHandlerFor(selectedRange))
		}(i, &r)
	}

	statusServeMux := http.NewServeMux()
	statusServeMux.HandleFunc("/", statusHandler)
	http.ListenAndServe("localhost:8080", statusServeMux)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := Status{
		Status:      "okay",
		Advertizing: "absolutely",
	}
	if err := json.NewEncoder(w).Encode(status); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func setupHandlerFor(selectedRange *RangeDetails) *http.ServeMux {
	resultServeMux := http.NewServeMux()
	resultServeMux.HandleFunc("/", getRangeAdvertHandler(selectedRange))

	return resultServeMux
}

func getRangeAdvertHandler(selectedRange *RangeDetails) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(&selectedRange); err == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
