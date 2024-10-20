package handler

import (
	"os"
	"strconv"

	"github.com/unickorn/footcal/sofa"
	"github.com/appwrite/sdk-for-go/databases"
)

// DB is a wrapper around the appwrite database client.
type DB struct {
	db *databases.Databases
}

// NewDB creates a new DB from an appwrite database client.
func NewDB(dbs *databases.Databases) *DB {
	return &DB{
		db: dbs,
	}
}

// FindMatch returns a match from the appwrite database.
func (db *DB) FindMatch(id int64) (sofa.Match, bool) {
	doc, err := db.db.GetDocument(os.Getenv("APPWRITE_DB_ID"), "matches", strconv.Itoa(int(id)))
	if err != nil {
		return sofa.Match{}, false
	}

	var match sofa.Match
	err = doc.Decode(&match)
	if err != nil {
		return sofa.Match{}, false
	}

	return match, true
}

// SaveMatch saves a match to the appwrite database.
func (db *DB) SaveMatch(m sofa.Match) error {
	_, err := db.db.CreateDocument(os.Getenv("APPWRITE_DB_ID"), "matches", strconv.Itoa(int(m.ID)), m)
	return err
}
