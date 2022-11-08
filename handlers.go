package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	status := Status{
		Status:      "okay",
		Advertising: "absolutely",
	}
	if err := json.NewEncoder(w).Encode(status); err == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getRangeAdvertJsonHandler(selectedRange *RangeDetails) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(&selectedRange); err == nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func getRangeAdvertUiHandler(selectedRange *RangeDetails) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		responseString, htmlGenErr := getHtmlForRange(selectedRange)
		if htmlGenErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_, err := io.WriteString(w, responseString)

		if err != nil {
			fmt.Println(err)
		}
	}
}

func getRangeAdvertImageHandler(selectedRange *RangeDetails) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		imageBytes, imageReadErr := assets.ReadFile(fmt.Sprintf("assets/%s.png", selectedRange.Model))
		if imageReadErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Write(imageBytes)
		}
	}
}

func getRangeAdvertFaviconHandler(selectedRange *RangeDetails) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(favicon)
	}
}

func logRequestHandlerWrapper(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		fmt.Println(generateLogLine(r.URL.String(), r.Method, r.Proto, r.RemoteAddr))
	}

	return http.HandlerFunc(fn)
}
