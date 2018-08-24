package models

import (
	"github.com/udemy_fileserver/datastore"
)

// File exported type structure with tags for database entries
type File struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}


type files struct{}

// Files exported type struct
var Files = new(files)

func (files) Create(name string) (*File, error) {
	f := &File{Name: name}
	var insertedID int
	tx := datastore.Postgre.MustBegin()
	tx.QueryRow("INSERT INTO files (name) VALUES ($1) returning id;", name).Scan(&insertedID)
	err := tx.Commit()
	f.ID = insertedID
	return f, err
}

func (files) List() ([]*File, error) {
	f := []*File{}
	err := datastore.Postgre.Select(&f, "SELECT id, name FROM files")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (files) ByName(name string) (*File, error) {
	var filetype File
	f := &filetype
	err := datastore.Postgre.Get(f, "SELECT id, name FROM files WHERE name=$1;", name)
	if err != nil {
		return nil, err
	}
	return f, nil
}
