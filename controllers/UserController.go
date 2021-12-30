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
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {

	jsonBodyMap := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&jsonBodyMap)

	exitRoute := return400OnError(w, err, "error processing request")

	if exitRoute {
		return
	}

	username, usernameNotNil := mapValueNotNil(jsonBodyMap, "email")
	password, passwordNotNil := mapValueNotNil(jsonBodyMap, "password")

	// validate username and password
	if usernameNotNil && passwordNotNil && interfaceIsString(username) && interfaceIsString(password) {

		user := models.User{}
		Db.First(&user, "email = ? AND password = ?", username, password)

		// user was found and is valid
		if user.Id > 0 && user.Status == "enabled" {

			// create the session object
			session := models.Session{}

			// delete any previous sessions for this user
			Db.Where("user_id = ?", user.Id).Delete(&session)

			// generate a new session ID
			b := make([]byte, 32)

			if _, err := io.ReadFull(rand.Reader, b); err != nil {
				w.WriteHeader(400)
				fmt.Fprintln(w, "Could not create session")
				return
			}

			// populate the session object
			session.ExpirationTime = time.Now().Add(time.Minute * time.Duration(20))
			session.Id = base64.URLEncoding.EncodeToString(b)
			session.UserId = user.Id

			// create the session in the database
			Db.Create(&session)

			// output the new session
			json.NewEncoder(w).Encode(session)
			return

		}

	}

	w.WriteHeader(400)
	fmt.Fprintln(w, "could not log in with given credentials")

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)

	var user models.User

	json.Unmarshal(reqBody, &user)
	Db.Create(&user)

	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userid"]

	users := []models.User{}
	Db.Find(&users)

	if len(userId) > 0 {

		for _, user := range users {
			// string to int
			userIdNum, err := strconv.Atoi(userId)
			if err == nil {
				if user.Id == userIdNum {
					json.NewEncoder(w).Encode(user)
				}
			}
		}

	} else {

		json.NewEncoder(w).Encode(users)

	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userid"]

	submittedUserId, conversionErr := strconv.Atoi(userId)

	// ensure that the id passed is a valid int
	if conversionErr != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Invalid user ID")
		return
	}

	// look up the user that we are updating
	user := models.User{}
	Db.First(&user, "id = ?", userId)

	var updatedUser models.User

	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &updatedUser)
	updatedUser.Id = submittedUserId

	if user.Id == updatedUser.Id {

		Db.Save(&updatedUser)

	} else {

		w.WriteHeader(400)
		fmt.Fprintln(w, "Error updating user ")
		return

	}

	json.NewEncoder(w).Encode(updatedUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userid"]

	_, conversionErr := strconv.Atoi(userId)

	if conversionErr != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, "Invalid user ID")
		return
	}

	user := models.User{}
	Db.First(&user, "id = ?", userId)

	user.Password = ""
	user.Email = ""

	if user.Id > 0 {

		Db.Delete(&user)
		json.NewEncoder(w).Encode(user)

	} else {

		w.WriteHeader(400)
		fmt.Fprintln(w, "Error deleting user")

	}

}
