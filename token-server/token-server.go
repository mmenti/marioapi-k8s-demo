package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type ApiToken struct {
	Token []struct {
		UserName  string `json:"user_name"`
		UserToken string `json:"user_token"`
	} `json:"token"`
}

func GetTokenInfo(token string) string {

	var token_data ApiToken

	file, e := ioutil.ReadFile("/tokens.json")
	if e != nil {
		return ""
	}

	err := json.Unmarshal([]byte(string(file)), &token_data)
	if err != nil {
		return ""
	}

	for _, tk := range token_data.Token {
		if tk.UserToken == token {
			return tk.UserName
		}
	}
	return ""
}

func serve(w http.ResponseWriter, r *http.Request) {
	// return token username for token, or empty string if token value not found
	fmt.Fprintf(w, GetTokenInfo(r.FormValue("token")))
}

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
