package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

type authenticationPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *server) authentication(w http.ResponseWriter, r *http.Request) {
	var reqPayload authenticationPayload

	err := s.readJSON(w, r, &reqPayload)
	if err != nil {
		s.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := s.UserDB.GetUserByEmail(reqPayload.Email)
	if err != nil {
		s.errorJSON(w, errors.New("couldn't find user in database"), http.StatusNotFound)
		return
	}

	err = s.UserDB.PasswordCheck(user.Password, reqPayload.Password)
	if err != nil {
		s.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	resPayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Signed in as %s!", user.Email),
		Data:    user,
	}

	err = s.writeJSON(w, resPayload, http.StatusOK)
	if err != nil {
		s.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	log.Println("authentication service: successful signin")
}
