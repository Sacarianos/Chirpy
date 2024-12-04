package main

import "net/http"

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Create a fileserver thet serves files from the current directory
	fs := http.FileServer(http.Dir("."))

	// Use the fileserver to handle requests to the root path
	mux.Handle("/", fs)

	// Create a new srrver struct
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux, //Use the ServeMux as the handler
	}

	//Start the server
	if err := server.ListenAndServe(); err != nil {
		panic(err) // Log if the server fails to start
	}
}
