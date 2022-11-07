package main

import (
	"encoding/json"
	"flag"
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

func getHtmlForRange(selectedRange *RangeDetails) (string, error) {
	htmlToReturn := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Document</title>
	<style>
		html {
			height: 100%;
		}

		body {
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol";
			background-color: #90caf9;
			height: 100%;
			display: flex;
			justify-content: center;
			align-items: center;
		}

		h1 {
			text-align: center;
		}
		
		a {
			color: #000;
		}
	</style>
</head>
<body>
  <h1>
		<span>Click here to buy the</span>
		<br />
		<a href="{{.Link}}" target="_blank">{{.Make}} {{.Model}}!</a>
	</h1>
</body>
</html>
`

	tmpl, parseErr := template.New("range").Parse(htmlToReturn)
	if parseErr != nil {
		return "", parseErr
	}
	buf := bytes.NewBuffer(make([]byte, 0))

	execErr := tmpl.Execute(buf, selectedRange)
	if execErr != nil {
		return "", execErr
	}

	return buf.String(), nil
}
