package app

import (
	"fmt"
	"net/http"
	"time"
)

// runs goku as a web server using the port argument
func Serve(port string) {

	// concats the address with the port number from the env variable
	address := "0.0.0.0:" + port

	// creates mux which is the router used for the http server
	mux := http.NewServeMux()

	// sets the endpoint /goku to run the function gokuHandler
	mux.HandleFunc("/goku", gokuHandler)

	// set the fields of the http server, the Handler will point to
	// the router
	server := http.Server{
		Addr:         address,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Goku is running on http://" + address + "/goku")

	// start the http server, it will stop execution if an error is returned
	panic(server.ListenAndServe())

}
