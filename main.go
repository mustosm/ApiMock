package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties"
)

type Mock struct {
	UUID    string       `json:"uuid,omitempty"`
	Message string       `json:"message,omitempty"`
	Headers []HttpHeader `json:"headers,omitempty"`
}

type HttpHeader struct {
	Value string `json:"value,omitempty"`
}

func GetMock(w http.ResponseWriter, r *http.Request) {
	var head []HttpHeader
	delay := r.URL.Query().Get("delay")
	i, err := strconv.Atoi(delay)
	if err == nil {
		time.Sleep(time.Duration(i) * time.Millisecond)
	}

	for key, value := range r.Header {
		head = append(head, HttpHeader{Value: strings.Join([]string{key, strings.Join(value, " ")}, ":")})
	}
	var mock Mock = Mock{UUID: uuid.New().String(), Message: "Hello world !", Headers: head}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mock)

}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
}

func main() {
	p := properties.MustLoadFile("ApiMock.properties", properties.UTF8)
	r := mux.NewRouter()

	GetMockHandler := http.HandlerFunc(GetMock)
	r.Handle("/mock", GetMockHandler).Methods("GET")
	r.HandleFunc("/status", GetStatus).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+p.MustGetString("port"), handlers.LoggingHandler(os.Stdout, r)))
	return
}
