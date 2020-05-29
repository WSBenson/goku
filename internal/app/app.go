package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type testRequest struct {
	Name string `json:"name"`
}

// struct to hold the new concatination of the user's name
type testResponse struct {
	Response string
}

func Serve() {

	// Set the default port to 3000
	portStr := "3000"

	// Bind the viper key "port" to the env variable "PORT" from the
	// dockerfile then get the value from the "PORT" env variable as
	// as it is not nil.
	viper.BindEnv("port", "PORT")
	if viper.Get("port") != nil {
		port := viper.Get("port")
		// Convert "PORT" env variable to a string
		portStr = fmt.Sprintf("%v", port)
	}

	fmt.Println("Goku is running on port " + portStr)

	// concats the address with the port number from the env variable
	address := "0.0.0.0:" + portStr

	// creates mux which is the router used for the http server
	mux := http.NewServeMux()
	// sets the endpoint /test to run the function testHandler
	mux.HandleFunc("/test", testHandler)

	// set the fields of the http server, the Handler will point to
	// the router
	goSrvr := http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// start the http server, it will stop execution if an error is returned
	panic(goSrvr.ListenAndServe())

}

// testHandler accepts a post request that takes the user's name from the JSON request
// and responds back with "Hello" + that user's name
func testHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Reads the body of the JSON request into the testRequest struct
		// by unmarshaling the JSON bytes
		tr := testRequest{}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(body, &tr)

		// Concatinates the name from the POST body with "Hello"
		// then marshals the response back into JSON bytes to
		// send to the browser.
		tres := testResponse{
			Response: "Hello, " + tr.Name,
		}
		b, err := json.Marshal(tres)
		if err != nil {
			log.Fatal(err)
		}

		// Writes the JSON bytes to the JSON response writer
		w.Write(b)

	case http.MethodGet:
		fmt.Fprintln(w, "Get not supported")
	}

}
