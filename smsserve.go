/*
 *  Serves REST endpoints to store and receive SMS message from the android app:
 *  "SMS an PC/Telefon - Automatische Umleitung"
 *
 * GET /last - show the last received message as plain text only once
 * TODO: list - json of all received message since server restart
 * TODO: user auth and multiple user support
 * POST /store - store a message json: {"subject": "testsubject", "message": "you are not alone"}
 *                - Jessica Sommer
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type sms struct {
	UserName   string    `json:"user-name"`
	Message    string    `json:"message"`
	Subject    string    `json:"subject"`
	StoredDate time.Time `json:"stored-date"`
}

type user struct {
	UserName string
	Password string
}

var messages []sms
var users []user

func handleSms(w http.ResponseWriter, r *http.Request) {
	//Allow CORS here By * or specific origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	switch r.Method {
	case "GET": // write last received message to response and then clear stored message
		var auth, userName, authErr = paramterAuth(r.URL.Query())
		if authErr != nil || auth == false {
			http.Error(w, authErr.Error(), http.StatusForbidden)
			return
		}
		var message = ""
		//Find previously stored SMS for this user
		for _, sms := range messages {
			if sms.UserName == userName {
				message = sms.Message
			}
		}
		fmt.Fprintf(w, message)
	case "POST": //receive sms message and store text
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println(" Recived message : %s", string(body))
		var s sms
		err = json.Unmarshal(body, &s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var auth, userName, authErr = paramterAuth(r.URL.Query())
		if authErr != nil || auth == false {
			http.Error(w, authErr.Error(), http.StatusForbidden)
			return
		}

		//Find previously stored SMS for this user
		for i, sms := range messages {
			if sms.UserName == userName {
				//overwrite stored with new message
				messages[i].Message = s.Message
				return
			}
		}

		//No stored one found create one
		messages = append(messages, sms{UserName: userName, Message: s.Message, StoredDate: time.Now()})
		log.Println(s.Message)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

//Provide authentication via parameters as the "SMS an PC/Telefon"
//app does not support any authentication like Basic Auth, etc
func paramterAuth(values url.Values) (bool, string, error) {
	var pass = ""
	var userName = ""
	for k, v := range values {
		fmt.Println(k, " => ", v)
		if k == "user" && len(v) > 0 {
			userName = v[0]
		}
		if k == "pass" && len(v) > 0 {
			pass = v[0]
		}
	}
	if userName == "" || pass == "" {
		return false, "", errors.New("missing user or password")
	}

	for _, user := range users {
		if user.UserName == userName {
			if user.Password == pass {
				return true, userName, nil
			} else {
				return false, "", errors.New("wrong password")
			}
		}
	}
	return false, "", errors.New("unknown user")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadUsersFromFile() {
	f, err := ioutil.ReadFile("users")
	check(err)

	lines := strings.Split(string(f), "\n")
	fmt.Print(string(f))
	for _, line := range lines {
		userData := strings.Split(line, ":")
		if len(userData) > 1 {
			users = append(users, user{UserName: userData[0], Password: userData[1]})
			fmt.Printf("Load user: %s with passsword: %s\n", userData[0], userData[1])
		}
	}
}

func handleRequests() {
	log.Println("handling Request")
	http.HandleFunc("/messages", handleSms)
	http.HandleFunc("/messages/", handleSms)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func main() {
	loadUsersFromFile()
	handleRequests()
}
