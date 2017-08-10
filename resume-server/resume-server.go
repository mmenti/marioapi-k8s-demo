package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func LoadResume() string {

	file, e := ioutil.ReadFile("/resume.json")
	if e != nil {
		return "{}"
	}
	return string(file)
}

func serve(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, LoadResume())
}

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
