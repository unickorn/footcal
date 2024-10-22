package handler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/appwrite/sdk-for-go/databases"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
	"github.com/unickorn/footcal/sofa"
)

// DB is a wrapper around the appwrite database client.
type DB struct {
	db *databases.Databases
	ctx openruntimes.Context
}

// NewDB creates a new DB from an appwrite database client.
func NewDB(dbs *databases.Databases, ctx openruntimes.Context) *DB {
	return &DB{
		db: dbs,
		ctx: ctx,
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
		db.ctx.Error(fmt.Sprintf("Failed to decode match: %s\n", err.Error()))
		return sofa.Match{}, false
	}

	db.ctx.Log(fmt.Sprintf("Found match: %v\n", match))
	return match, true
}

// SaveMatch saves a match to the appwrite database.
func (db *DB) SaveMatch(m sofa.Match) error {
	db.ctx.Log(fmt.Sprintf("Saving to DB ID:", os.Getenv("APPWRITE_DB_ID")))

	saved, err := db.db.CreateDocument(os.Getenv("APPWRITE_DB_ID"), "matches", strconv.Itoa(int(m.ID)), m)
	db.ctx.Log(fmt.Sprintf("Saved match: %v -> %v\n", m, saved.Id))
	return err
}
