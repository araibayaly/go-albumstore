Album Store

Introduction
The Go Album Store project is a simple web application for managing albums. It provides CRUD (Create, Read, Update, Delete) operations for albums, a RESTful API using Gorilla Mux, and PostgreSQL as the backend database.

API Documentation

The Go Album Store API supports the following endpoints:

GET /albums: Retrieve all albums.
GET /albums/{id}: Retrieve a specific album by ID.
POST /albums: Create a new album.
PUT /albums/{id}: Update an existing album by ID.
DELETE /albums/{id}: Delete an album by ID.

DB Structure
Table restaurants {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  title text
  artist text
  genre text
  year text
}
