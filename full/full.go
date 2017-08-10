package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var ()

func LoadResume() (resumeData Resume, loadedOk bool) {

	var res Resume

	rsp, err := http.Get("http://resumeserver")
	if err != nil {
		return res, false
	}
	defer rsp.Body.Close()
	bodyByte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return res, false
	}

	err = json.Unmarshal(bodyByte, &resumeData)
	if err != nil {
		return res, false
	}
	return resumeData, true

}

func WriteApiError(w http.ResponseWriter, errorCode int, errorStr string) {

	apierror := ApiError{errorCode, errorStr}
	json, err := json.Marshal(apierror)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func WriteApiSuccess(w http.ResponseWriter, successCode int, successStr string) {

	apisuccess := ApiSuccess{successCode, successStr}
	json, err := json.Marshal(apisuccess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func serve(w http.ResponseWriter, r *http.Request) {

	// check token, load resume, return relevant part(s)
	token := r.FormValue("token")
	rsp, err := http.Get("http://tokenserver?token=" + token)
	if err != nil {
		WriteApiError(w, 110, "Error checking token from token server")
		return
	}
	defer rsp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		WriteApiError(w, 111, "Error checking token from token server")
		return
	}
	if string(bodyBytes) == "" {
		WriteApiError(w, 401, "Invalid token "+token)
		return
	}

	resumeData, loadedOk := LoadResume()

	if !loadedOk {
		WriteApiError(w, 101, "Error loading resume")
		return
	} else {
		json, err := json.Marshal(resumeData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
