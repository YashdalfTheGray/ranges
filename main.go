package main

import (
	_ "embed"
	"flag"
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

//go:embed range.tpl.html
var htmlTemplate string

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
	var bindAddr string

	flag.StringVar(&bindAddr, "bind-address", "localhost", "the address to bind the ports to")
	flag.Parse()

	for i, r := range allRanges {
		i := i
		r := r
		go func(port int, selectedRange *RangeDetails) {
			http.ListenAndServe(fmt.Sprintf("%s:808%d", bindAddr, port), setupHandlerFor(selectedRange))
		}(i+1, &r)
	}
	fmt.Println("Started advert servers")

	statusServeMux := http.NewServeMux()
	statusServeMux.HandleFunc("/", statusHandler)
	wrappedStatusHandler := logRequestHandlerWrapper(statusServeMux)
	fmt.Println("Started status server")
	http.ListenAndServe(fmt.Sprintf("%s:8080", bindAddr), wrappedStatusHandler)
}

func setupHandlerFor(selectedRange *RangeDetails) http.Handler {
	resultServeMux := http.NewServeMux()
	resultServeMux.HandleFunc("/json", getRangeAdvertJsonHandler(selectedRange))
	resultServeMux.HandleFunc("/", getRangeAdvertUiHandler(selectedRange))

	wrappedHandler := logRequestHandlerWrapper(resultServeMux)

	return wrappedHandler
}
