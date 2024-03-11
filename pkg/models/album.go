package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Album struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Title     string `json:"name"`
	Artist    string `json:"artist"`
	Genre     string `json:"genre"`
	Year      string `json:"year"`
}

type AlbumModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func (a AlbumModel) Insert(album *Album) error {
	query := `
		INSERT INTO albums (title, artist, genre, year) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id, created_at, updated_at
		`
	args := []interface{}{album.Title, album.Artist, album.Genre, album.Year}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&album.Id, &album.CreatedAt, &album.UpdatedAt)
}

func (a AlbumModel) Get(id int) (*Album, error) {
	query := `
		SELECT id, created_at, updated_at, title, artist, genre, year
		FROM albums
		WHERE id = $1
		`
	var album Album
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := a.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&album.Id, &album.CreatedAt, &album.UpdatedAt, &album.Title, &album.Artist, &album.Genre, &album.Year)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (a AlbumModel) GetAll() ([]*Album, error) {
	query := `
		SELECT id, created_at, updated_at, title, artist, genre, year
		FROM albums
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := a.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var albums []*Album

	for rows.Next() {
		var album Album
		err := rows.Scan(&album.Id, &album.CreatedAt, &album.UpdatedAt, &album.Title, &album.Artist, &album.Genre, &album.Year)
		if err != nil {
			return nil, err
		}
		albums = append(albums, &album)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return albums, nil
}

func (a AlbumModel) Update(album *Album) error {
	query := `
		UPDATE albums
		SET title = $1, artist = $2, genre = $3, year = $4 
		WHERE id = $5
		RETURNING updated_at
		`
	args := []interface{}{album.Title, album.Artist, album.Genre, album.Year, album.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return a.DB.QueryRowContext(ctx, query, args...).Scan(&album.UpdatedAt)
}

func (a AlbumModel) Delete(id int) error {
	query := `
		DELETE FROM albums
		WHERE id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := a.DB.ExecContext(ctx, query, id)
	return err
}
