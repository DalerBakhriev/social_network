package apiserver

import (
	"database/sql"
	"net/http"

	"github.com/DalerBakhriev/social_network/internal/app/store/sqlstore"
	"github.com/gorilla/sessions"
)

func newDB(databaseURL string) (*sql.DB, error) {

	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Start ...
func Start(config *Config) error {

	db, err := newDB(config.DataBaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}
