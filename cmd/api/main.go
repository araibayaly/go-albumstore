package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	model "github.com/araibayaly/go-albumstore/pkg/models"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

type config struct {
	port string
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	models model.Model
}

func main() {
	var cfg config
	flag.StringVar(&cfg.port, "port", ":8081", "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://postgres:1997@localhost/albumstore_db?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Connect to DB
	db, err := openDB(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	app := &application{
		config: cfg,
		models: model.NewModels(db),
	}

	app.run()
}

func (app *application) run() {
	r := mux.NewRouter()

	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/albums", app.createAlbumHandler).Methods("POST")
	v1.HandleFunc("/albums/{id}", app.getAlbumHandler).Methods("GET")
	v1.HandleFunc("/albums", app.getAllAlbumsHandler).Methods("GET")
	v1.HandleFunc("/albums/{id}", app.updateAlbumHandler).Methods("PUT")
	v1.HandleFunc("/albums/{id}", app.deleteAlbumHandler).Methods("DELETE")

	log.Printf("Starting server on %s\n", app.config.port)
	err := http.ListenAndServe(app.config.port, r)
	log.Fatal(err)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
