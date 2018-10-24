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
	encoderr := json.NewEncoder(w).Encode(mock)
	if encoderr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w)
	}

}

func GetStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is up and running"))
}

func GetSwagger (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w,r,"apimock-swagger.json")
}

func main() {
	p := properties.MustLoadFile("ApiMock.properties", properties.UTF8)
	r := mux.NewRouter()

	GetMockHandler := http.HandlerFunc(GetMock)
	r.Handle("/mock", GetMockHandler).Methods("GET")
	r.HandleFunc("/health", GetStatus).Methods("GET")
	r.HandleFunc("/swagger", GetSwagger).Methods("GET")
	log.Fatal(http.ListenAndServeTLS(":"+p.MustGetString("port"),p.MustGetString("certificate"),p.MustGetString("key"),handlers.LoggingHandler(os.Stdout, r)))
	return
}
