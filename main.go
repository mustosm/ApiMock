package main

import (
	"fmt"
	"strings"
    "encoding/json"
    "log"
	"net/http"
	"time"
	"strconv"
	"github.com/google/uuid"
	"github.com/magiconair/properties"
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
	if r.Method == "GET" {
		var head []HttpHeader
		delay := r.URL.Query().Get("delay")
		i, err := strconv.Atoi(delay)
		if err == nil {	
			time.Sleep(time.Duration(i) * time.Millisecond)
		}
		fields := strings.Split(r.URL.Query().Get("fields"), ",")
		fmt.Println(len(fields))
		for key, value := range r.Header {		
			head = append(head, HttpHeader{Value: strings.Join([]string{key, strings.Join(value," ")}, ":")})
		}
		var mock Mock
		if len(fields) > 0 {
			for _,f := range fields {
				if (f == "UUID") {
					mock.UUID = uuid.New().String()
				}
				if (f == "Message") {
					mock.Message = "Hello world !"
				}
				if (f == "Headers") {
					mock.Headers = head
				} 
			}
		} else {
			mock = Mock{UUID: uuid.New().String(), Message: "Hello world !", Headers: head}
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)	
		json.NewEncoder(w).Encode(mock)
	} else {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)		
	}
}

func main() {
	p := properties.MustLoadFile("ApiMock.properties", properties.UTF8)
	mux := http.NewServeMux()
	mux.HandleFunc("/mock", GetMock)
	log.Fatal(http.ListenAndServe(":"+p.MustGetString("port"),mux))
	return

}