package apiserver

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"io/ioutil"

	"github.com/DalerBakhriev/social_network/internal/app/model"
	"github.com/gorilla/mux"
)

const (
	numUsersOnOnePage = 20
	loginFormFilePath = "./templates/login.html"
)

func (s *server) handleSignUp() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			loginForm, err := ioutil.ReadFile(loginFormFilePath)
			if err != nil {
				s.respond(w, r, http.StatusInternalServerError, err)
				return
			}
			w.Write(loginForm)
			return
		}

		inputEmail := r.FormValue("email")
		inputPassword := r.FormValue("password")

		user := &model.User{
			Email:    inputEmail,
			Password: inputPassword,
		}

		if err := s.store.User().Create(user); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		user.Sanitize()
		s.respond(w, r, http.StatusCreated, user)
	}
}

func (s *server) handleMainPage() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		users, err := s.store.User().GetTopUsers(numUsersOnOnePage)
		if err != nil {
			s.respond(w, r, http.StatusInternalServerError, err)
			return
		}
		tmpl := template.Must(template.ParseFiles("./templates/users.html"))
		tmpl.Execute(w, users)
	}

}

func (s *server) handleLogIn() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			loginForm, err := ioutil.ReadFile(loginFormFilePath)
			if err != nil {
				s.respond(w, r, http.StatusInternalServerError, err)
				return
			}
			w.Write(loginForm)
		}

		inputEmail := r.FormValue("email")
		inputPassword := r.FormValue("password")

		user, err := s.store.User().FindByEmail(inputEmail)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		if passwordIsCorrect := user.ComparePassword(inputPassword); !passwordIsCorrect {
			s.error(w, r, http.StatusUnauthorized, errInncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.New(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, errors.New("Failed to create new session"))
			return
		}

		session.Values["user_id"] = user.ID

		http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID), http.StatusPermanentRedirect)
	}
}

func (s *server) handleGetSingleUser() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("./templates/user.html"))
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		userID := vars["user_id"]
		id, err := strconv.Atoi(userID)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		user, err := s.store.User().Find(id)
		tmpl.Execute(w, user)
	}

}

func (s *server) handleGetFriendsList() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles("./templates/user.html"))
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		userID := vars["user_id"]

		id, err := strconv.Atoi(userID)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		users, err := s.store.User().GetFriendsList(id)
		tmpl.Execute(w, users)
	}
}

func (s *server) getUserID(w http.ResponseWriter, r *http.Request) (int, error) {

	session, err := s.sessionStore.Get(r, sessionName)

	if err != nil {
		return -1, err
	}

	userID, ok := session.Values["user_id"]
	if !ok {
		return -1, errNotAuthenticated
	}

	return userID.(int), nil
}

func (s *server) getFriendID(w http.ResponseWriter, r *http.Request) (int, error) {

	friendID := r.FormValue("friend_id")
	if friendID == "" {
		return -1, errors.New("Could not find friend id")
	}

	idTo, err := strconv.Atoi(friendID)
	if err != nil {
		return -1, err
	}

	return idTo, nil
}

func (s *server) handleSendFriendsRequest() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := s.getUserID(w, r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		friendID, err := s.getFriendID(w, r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		errFriendRequest := s.store.User().SendFriendRequest(userID, friendID)
		if errFriendRequest != nil {
			s.error(w, r, http.StatusInternalServerError, errFriendRequest)
			return
		}

	}
}

func (s *server) handleAcceptFriendsRequest() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := s.getUserID(w, r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		friendID, err := s.getFriendID(w, r)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.User().AcceptFriendRequest(userID, friendID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/users/%d", userID), http.StatusFound)
	}
}
