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

func serve(w http.ResponseWriter, r *http.Request) {

	urlArg := r.URL.Path[len("/"):]

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
		if len(urlArg) <= 0 {
			tgs := Tags{}
			// go through all projects and experience and collect tags
			for i, proj := range resumeData.Projects {
			TAGLOOP1:
				for j, tags := range proj.Tags {
					if i == 0 && j == 0 {
						tg := &Tag{}
						tg.Name = tags
						tg.Count = 1
						tgs = append(tgs, tg)
					} else {
						for _, extags := range tgs {
							if extags.Name == tags {
								extags.Count++
								continue TAGLOOP1
							}
						}
						tg := &Tag{}
						tg.Name = tags
						tg.Count = 1
						tgs = append(tgs, tg)
					}
				}
			}
			for i, proj := range resumeData.Experience {
			TAGLOOP2:
				for j, tags := range proj.Tags {
					if i == 0 && j == 0 {
						tg := &Tag{}
						tg.Name = tags
						tg.Count = 1
						tgs = append(tgs, tg)
					} else {
						for _, extags := range tgs {
							if extags.Name == tags {
								extags.Count++
								continue TAGLOOP2
							}
						}
						tg := &Tag{}
						tg.Name = tags
						tg.Count = 1
						tgs = append(tgs, tg)
					}
				}
			}

			json, err := json.Marshal(tgs)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)

		} else {
			// find projects and experience that matches this tag
			sr := SearchResult{}
			for _, proj := range resumeData.Projects {
				for _, tags := range proj.Tags {
					if urlArg == tags {
						sr.Projects = append(sr.Projects, proj)
					}
				}
			}
			for _, proj := range resumeData.Experience {
				for _, tags := range proj.Tags {
					if urlArg == tags {
						sr.Experience = append(sr.Experience, proj)
					}
				}
			}
			json, err := json.Marshal(sr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)

		}
	}
}

func main() {
	http.HandleFunc("/", serve)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
