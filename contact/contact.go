package main

import (
	"encoding/json"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	// twilio (for SMS) and SendGrid (for email) config
	// passed via a Secret to the env
	twilioSid    string = os.Getenv("TWILIO_SID")
	twilioToken  string = os.Getenv("TWILIO_TOKEN")
	twilioUrl    string = os.Getenv("TWILIO_URL")
	twilioNumber string = os.Getenv("TWILIO_NUMBER")
	alertNumber  string = os.Getenv("ALERT_NUMBER")
	sendGridKey  string = os.Getenv("SENDGRID_KEY")
)

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

func SendSMS(w http.ResponseWriter, messageTxt string) {

	v := url.Values{}
	v.Set("To", alertNumber)
	v.Set("From", twilioNumber)
	v.Set("Body", messageTxt)
	rb := *strings.NewReader(v.Encode())

	client := &http.Client{}

	req, e := http.NewRequest("POST", twilioUrl, &rb)
	if e != nil {
		WriteApiError(w, 110, "I'm sorry, there was problem with sending your SMS : "+e.Error())
		return
	}
	req.SetBasicAuth(twilioSid, twilioToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		WriteApiError(w, 110, "I'm sorry, there was problem with sending your SMS : "+err.Error())
		return
	}

	if resp.Status == "200" {
		WriteApiError(w, 110, "I'm sorry, there was problem with sending your SMS")
	} else {
		WriteApiSuccess(w, resp.StatusCode, "SMS successfully sent, thanks so much!")
	}
}

func SendEmail(w http.ResponseWriter, messageTxt string, fromAddr string) {

	from := mail.NewEmail("Mario's API", fromAddr)
	subject := "Email message from api.mariomenti.com"
	to := mail.NewEmail("Mario Menti", "mario@menti.net")
	plainTextContent := messageTxt
	htmlContent := "<p>" + messageTxt + "</p>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(sendGridKey)
	response, err := client.Send(message)
	if err != nil {
		WriteApiError(w, 111, fmt.Sprintf("I'm sorry, there was problem with sending your email via SendGrid: %v", err))

	} else {
		WriteApiSuccess(w, response.StatusCode, "Email successfully sent, thanks so much!")
	}

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

	var ct Contact
	resumeData, loadedOk := LoadResume()

	if !loadedOk {
		WriteApiError(w, 101, "Error loading resume")
		return
	} else {
		switch r.Method {
		case "GET":
			ct.Name = resumeData.Name
			ct.ContactPhone = resumeData.ContactPhone
			ct.ContactEmail = resumeData.ContactEmail
			ct.ContactAddress = resumeData.ContactAddress
			json, err := json.Marshal(ct)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)

		case "POST":
			channel := r.FormValue("channel")
			message := r.FormValue("message")
			from := r.FormValue("from")

			if channel != "sms" && channel != "email" {
				WriteApiError(w, 102, "Parameter 'channel' needs to be one of 'sms' or 'email'")
				return
			}
			if channel == "email" && from == "" {
				WriteApiError(w, 102, "When specifying the email 'channel', please also provide the 'from' parameter so Mario can reply to you :)")
				return
			}
			if message == "" {
				WriteApiError(w, 102, "This endpoint requires a 'message' parameter")
				return
			}

			if channel == "sms" {
				SendSMS(w, message)
			} else {
				if channel == "email" {
					SendEmail(w, message, from)
				}
			}
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
