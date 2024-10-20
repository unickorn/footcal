package main

import (
	"github.com/unickorn/footcal/sofa"
)

// DB is a wrapper around the appwrite database client.
type DB struct {
}

// FindMatch returns a match from the appwrite database.
func (db DB) FindMatch(id int64) (sofa.Match, bool) {
	return sofa.Match{}, false
}

// SaveMatch saves a match to the appwrite database.
func (db DB) SaveMatch(m sofa.Match) error {
	return nil
}
