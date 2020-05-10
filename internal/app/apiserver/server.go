package apiserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/DalerBakhriev/social_network/internal/app/store"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

const (
	sessionName        = "some_session"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *zap.SugaredLogger
	store        store.Store
	sessionStore sessions.Store
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(store store.Store, sessionStore sessions.Store) *server {

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Can't initializw zap logger %v", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	s := &server{
		router:       mux.NewRouter(),
		logger:       sugaredLogger,
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) error(w http.ResponseWriter, r *http.Request, statusCode int, err error) {
	s.respond(w, r, statusCode, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) configureRouter() {

	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.HandleFunc("/signup", s.handleSignUp()).Methods("POST")
	s.router.HandleFunc("/login", s.handleLogIn()).Methods("POST")
	s.router.HandleFunc("/", s.handleMainPage()).Methods("GET")
	s.router.HandleFunc("/users/{user_id:[0-9]+}", s.handleGetSingleUser()).Methods("GET")
	s.router.HandleFunc("/users/{user_id:[0-9]+}/friend_request", s.handleSendFriendsRequest()).Methods("POST")
	s.router.HandleFunc("users/{user_id:[0-9]+}/accept_friend_request", s.handleAcceptFriendsRequest()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
}
