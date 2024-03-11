package model

import (
	"database/sql"
	"log"
	"os"
)

type Model struct {
	Albums AlbumModel
}

func NewModels(db *sql.DB) Model {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Model{
		Albums: AlbumModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
