package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	model "github.com/araibayaly/go-albumstore/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string `json:"title"`
		Artist string `json:"artist"`
		Genre  string `json:"genre"`
		Year   string `json:"year"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	album := &model.Album{
		Title:  input.Title,
		Artist: input.Artist,
		Genre:  input.Genre,
		Year:   input.Year,
	}

	err = app.models.Albums.Insert(album)
	if err != nil {
		log.Println("Error inserting album:", err)
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, album)
}

func (app *application) getAlbumHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		log.Println("Error inserting album:", err)
		app.respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	album, err := app.models.Albums.Get(id)
	if err != nil {
		if err == sql.ErrNoRows {
			app.respondWithError(w, http.StatusNotFound, "Album not found")
			return
		}
		// Handle other database errors
		log.Println("Error fetching album:", err)
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, album)
}

func (app *application) getAllAlbumsHandler(w http.ResponseWriter, r *http.Request) {
	albums, err := app.models.Albums.GetAll()
	if err != nil {
		// Handle database errors
		log.Println("Error fetching albums:", err)
		app.respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, albums)
}

func (app *application) updateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	album, err := app.models.Albums.Get(id)
	if err != nil {
		log.Println("Error inserting album:", err)
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Title  *string `json:"title"`
		Artist *string `json:"artist"`
		Genre  *string `json:"genre"`
		Year   *string `json:"year"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		album.Title = *input.Title
	}

	if input.Artist != nil {
		album.Artist = *input.Artist
	}

	if input.Genre != nil {
		album.Genre = *input.Genre
	}

	if input.Year != nil {
		album.Year = *input.Year
	}

	err = app.models.Albums.Update(album)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, album)
}

func (app *application) deleteAlbumHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = app.models.Albums.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
