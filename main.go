package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
		i := i
		r := r
		go func(port int, selectedRange *RangeDetails) {
			http.ListenAndServe(fmt.Sprintf("localhost:808%d", port), setupHandlerFor(selectedRange))
		}(i+1, &r)
	}
	fmt.Println("Started advert servers")

	statusServeMux := http.NewServeMux()
	statusServeMux.HandleFunc("/", statusHandler)
	fmt.Println("Started status server")
	http.ListenAndServe("localhost:8080", statusServeMux)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := Status{
		Status:      "okay",
		Advertizing: "absolutely",
	}
	if err := json.NewEncoder(w).Encode(status); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func setupHandlerFor(selectedRange *RangeDetails) http.Handler {
	resultServeMux := http.NewServeMux()
	resultServeMux.HandleFunc("/", getRangeAdvertHandler(selectedRange))

	wrappedHandler := logRequestHandlerWrapper(resultServeMux)

	return wrappedHandler
}

func getRangeAdvertHandler(selectedRange *RangeDetails) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(&selectedRange); err == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func logRequestHandlerWrapper(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		fmt.Println(generateLogLine(r.URL.String(), r.Method, r.Proto, r.RemoteAddr))
	}

	return http.HandlerFunc(fn)
}

func generateLogLine(uri, method, protocol, remote string) string {
	return fmt.Sprintf("[%s] \"%s %s %s\" %s", time.Now().UTC(), method, uri, protocol, remote)
}
