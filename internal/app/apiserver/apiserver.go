package apiserver

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

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

func getEnvOrDefaultValue(envVar, defaultVal string) string {

	envValue, ok := os.LookupEnv(envVar)
	if !ok {
		envValue = defaultVal
	}

	return envValue
}

func getDataBaseURL() string {

	dbHost := getEnvOrDefaultValue("MYSQL_HOST", "db")
	dbUser := getEnvOrDefaultValue("MYSQL_USER", "user")
	dbPassword := getEnvOrDefaultValue("MYSQL_PASSWORD", "password")
	dbName := getEnvOrDefaultValue("MYSQL_DATABASE", "social_network_db")

	dataBaseURL := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)

	return dataBaseURL
}

// Start ...
func Start(config *Config) error {

	dataBaseURL := getDataBaseURL()
	db, err := newDB(dataBaseURL)
	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	srv := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}
