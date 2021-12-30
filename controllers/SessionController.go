package controllers

import (
	"ThePooReview/models"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

/*
 * Generates a session with a unique ID and attaches it to the User object
 */
func CreateSession(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)

	var session models.Session

	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Could not create session")
		return
	}

	json.Unmarshal(reqBody, &session)
	session.ExpirationTime = time.Now().Add(time.Minute * time.Duration(20))
	session.Id = base64.URLEncoding.EncodeToString(b)

	Db.Create(&session)

	json.NewEncoder(w).Encode(session)
}

func GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userid"]

	session := models.Session{}
	Db.First(&session, "user_id=?", userId)

	if len(session.Id) > 0 {

		json.NewEncoder(w).Encode(session)

	} else {

		w.WriteHeader(400)
		fmt.Fprintln(w, "Invalid session")

	}
}
