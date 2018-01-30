package main

import (
	"strings"
    "encoding/json"
    "log"
    "net/http"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

type Mock struct {
    UUID    string   `json:"uuid,omitempty"`
	Message string   `json:"message,omitempty"`
	Headers []HttpHeader `json:"headers,omitempty"`
}

type HttpHeader struct {
	Value string `json:"value,omitempty"`
}

func GetMock(w http.ResponseWriter, r *http.Request) {
	var head []HttpHeader
	for key, value := range r.Header {		
		head = append(head, HttpHeader{Value: strings.Join([]string{key, strings.Join(value," ")}, ":")})
	}
	var mock Mock = Mock{UUID: uuid.New().String(), Message: "Hello world !", Headers: head}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)	
    json.NewEncoder(w).Encode(mock)
}

// our main function
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/mock", GetMock).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
	return
}