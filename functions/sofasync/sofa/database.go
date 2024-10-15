package sofa

// Database is an interface representing the database used alongside this package.
type Database interface {
	// FindMatch returns an event by its ID from the database, and a boolean indicating
	// whether it was found.
	FindMatch(id int64) (Match, bool)
	// SaveMatch saves a match onto the database.
	SaveMatch(m Match) error
}