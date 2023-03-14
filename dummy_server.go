package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

var db []string

func postData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handle /post request")
	if r.Method == "OPTIONS" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Origin", "*") //r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.WriteHeader(200)
		io.WriteString(w, "OK")
	}

	if r.Method != "POST" {
		w.WriteHeader(405)
		w.Header().Set("Allow", "POST")
		w.Write([]byte("Method not allowed"))
		fmt.Printf(" Invalid method, has %s instead of POST\n", r.Method)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("could not readt body: %s\n", err)
	}

	s := fmt.Sprintf("%s", body)

	fmt.Printf("Post body: %s\n", s)
	db = append(db, s)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*") //r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.WriteHeader(200)
	n := fmt.Sprintf("%s", rand.Int())
	io.WriteString(w, "{\"id\":\""+n+"\"}")
}

func getData(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Handle /get request\n")
	if r.Method != "GET" {
		w.WriteHeader(405)
		w.Header().Set("Allow", "GET")
		w.Write([]byte("Method not allowed"))
		fmt.Printf(" Invalid method, has %s instead of GET\n", r.Method)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*") //r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.WriteHeader(200)
	data := prepareData(db)
	fmt.Printf("Returning data: %s\n", data)
	io.WriteString(w, data)
}

func prepareData(data []string) string {
	s := "["
	for i := range data {
		s += data[i]
		s += ","
	}

	if len(s) > 1 {
		s = s[:len(s)-1]
	}
	s += "]"
	return s
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/get", getData)
	mux.HandleFunc("/post", postData)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server listening on", server.Addr)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
