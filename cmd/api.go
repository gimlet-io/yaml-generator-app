package main

import "net/http"

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func yamlGenerator(w http.ResponseWriter, r *http.Request) {
	// TODO helm template w onechart/onechart from body values, and send back the kubernetes yaml
	w.WriteHeader(200)
}
