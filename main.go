package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/auth0-community/auth0"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties"
	jose "gopkg.in/square/go-jose.v2"
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

func LoadPublicKey(data []byte) (interface{}, error) {
	input := data

	block, _ := pem.Decode(data)
	if block != nil {
		input = block.Bytes
	}

	// Try to load SubjectPublicKeyInfo
	pub, err0 := x509.ParsePKIXPublicKey(input)
	if err0 == nil {
		return pub, nil
	}

	cert, err1 := x509.ParseCertificate(input)
	if err1 == nil {
		return cert.PublicKey, nil
	}

	return nil, fmt.Errorf("square/go-jose: parse error, got '%s' and '%s'", err0, err1)
}

/* Set up a global string for our secret */
var mySigningKey = []byte("secret")
var validator *auth0.JWTValidator

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Creates a configuration with the Auth0 information
		data, err := ioutil.ReadFile("mustosm.pem")
		if err != nil {
			panic("Impossible to read key form disk")
		}

		secret, err := LoadPublicKey(data)
		if err != nil {
			panic("Invalid provided key")
		}
		audience := []string{"https://api.mustosm.io/v1/myFirstAPI"}
		secretProvider := auth0.NewKeyProvider(secret)
		configuration := auth0.NewConfiguration(secretProvider, audience, "https://mustosm.eu.auth0.com/", jose.RS256)
		validator = auth0.NewValidator(configuration)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func main() {
	p := properties.MustLoadFile("ApiMock.properties", properties.UTF8)
	r := mux.NewRouter()

	GetMockHandler := http.HandlerFunc(GetMock)
	r.Handle("/mock", GetMockHandler).Methods("GET")
	r.HandleFunc("/status", GetStatus).Methods("GET")
	//r.HandleFunc("/get-token", GetToken).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+p.MustGetString("port"), handlers.LoggingHandler(os.Stdout, r)))
	return
}
