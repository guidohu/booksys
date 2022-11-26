package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 80, "The port to listen on")
	flag.Parse()

	http.Handle("/api/v2/ping", http.HandlerFunc(getPing))

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

type pingResponse struct {
	Status string `json:"status"`
}

func getPing(w http.ResponseWriter, r *http.Request) {
	resp := pingResponse{
		Status: "ok",
	}
	fmt.Printf("Creating response: %+v\n", resp)

	j, err := json.Marshal(&resp)
	if err != nil {
		fmt.Println("Error creating JSON response")
		http.Error(w, "create pong response", http.StatusInternalServerError)
	}
	fmt.Printf("Json response: %+v", j)

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, bytes.NewReader(j))
}
