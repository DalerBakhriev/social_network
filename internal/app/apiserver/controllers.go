package apiserver

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/DalerBakhriev/social_network/internal/app/model"
	"github.com/gorilla/mux"
)

const (
	numUsersOnOnePage = 20
	templatesPath     = "./internal/app/apiserver/templates"
)

func (s *server) handleSignUp() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			loginForm, err := ioutil.ReadFile(path.Join(templatesPath, "signup.html"))
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			w.Write(loginForm)
			return
		}

		inputEmail := r.FormValue("email")
		inputPassword := r.FormValue("password")
		inputName := r.FormValue("name")
		inputSurname := r.FormValue("surname")
		inputCity := r.FormValue("city")
		inputAge := r.FormValue("age")
		inputSex := r.FormValue("sex")
		inputInterests := r.FormValue("interests")

		age, err := strconv.Atoi(inputAge)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		user := &model.User{
			Email:     inputEmail,
			Password:  inputPassword,
			Name:      inputName,
			Surname:   inputSurname,
			City:      inputCity,
			Age:       age,
			Sex:       inputSex,
			Interests: inputInterests,
		}

		if err := s.store.User().Create(user); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		user.Sanitize()
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func (s *server) handleMainPage() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		users, err := s.store.User().GetTopUsers(numUsersOnOnePage)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		usersForTemplate := model.Users{Users: users}
		tmpl := template.Must(template.ParseFiles(path.Join(templatesPath, "users.html")))
		tmpl.Execute(w, usersForTemplate)
	}

}

func (s *server) handleLogIn() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			loginForm, err := ioutil.ReadFile(path.Join(templatesPath, "login.html"))
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			w.Write(loginForm)
			return
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

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = user.ID
		err = session.Save(r, w)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID), http.StatusFound)
	}
}

func (s *server) handleUserEdit() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {

			userEditForm, err := ioutil.ReadFile(path.Join(templatesPath, "user_edit.html"))
			if err != nil {
				s.respond(w, r, http.StatusInternalServerError, err)
				return
			}
			w.Write(userEditForm)
			return
		}

		userID, err := s.getUserID(w, r)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		inputName := r.FormValue("name")
		inputSurname := r.FormValue("surname")
		inputCity := r.FormValue("city")
		inputAge := r.FormValue("age")
		inputSex := r.FormValue("sex")
		inputInterests := r.FormValue("interests")

		age, err := strconv.Atoi(inputAge)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errors.New("wrong age format, must be number"))
			return
		}

		user := &model.User{
			ID:        userID,
			Name:      inputName,
			Surname:   inputSurname,
			City:      inputCity,
			Age:       age,
			Sex:       inputSex,
			Interests: inputInterests,
		}

		if err := s.store.User().Update(user); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/users/%d", user.ID), http.StatusFound)
	}
}

func (s *server) handleGetSingleUser() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles(path.Join(templatesPath, "user.html")))
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		userID := vars["user_id"]
		id, err := strconv.Atoi(userID)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		user, err := s.store.User().Find(id)
		s.logger.Infof("Got user %v", user)
		tmpl.Execute(w, user)
	}

}

func (s *server) handleGetFriendsRequests() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles(path.Join(templatesPath, "requests.html")))
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		userID := vars["user_id"]

		id, err := strconv.Atoi(userID)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		users, err := s.store.User().GetFriendsList(id)

		usersForTemplate := model.FriendsAndRequests{
			Users:      users,
			CurrUserID: id,
		}
		tmpl.Execute(w, usersForTemplate)
	}
}

func (s *server) handleGetFriendsList() http.HandlerFunc {

	tmpl := template.Must(template.ParseFiles(path.Join(templatesPath, "friends.html")))
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		userID := vars["user_id"]

		id, err := strconv.Atoi(userID)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		users, err := s.store.User().GetFriendsList(id)

		usersForTemplate := model.FriendsAndRequests{
			Users:      users,
			CurrUserID: id,
		}
		tmpl.Execute(w, usersForTemplate)
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

	vars := mux.Vars(r)
	friendID := vars["friend_id"]

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

		vars := mux.Vars(r)
		userIDFromRequest, err := strconv.Atoi(vars["user_id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if userID != userIDFromRequest {
			s.error(w, r, http.StatusUnauthorized, errors.New("To accept this request login as current user"))
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
