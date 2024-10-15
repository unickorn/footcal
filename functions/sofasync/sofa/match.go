package sofa

import (
	"fmt"
	"time"
)

// Match is a struct containing information about a match.
type Match struct {
	// ID is the ID of the event in SofaScore.
	ID int64
	// Tournament is the name of the tournament this match finds place in.
	Tournament string
	// HomeTeam is the name of the home team.
	HomeTeam string
	// AwayTeam is the name of the away team.
	AwayTeam string
	// StartTime is the time the match will start.
	StartTime time.Time
	// Location is the location the match will find place in.
	Location string
}

// CollectMatches collects matches for a given team, using the db as a cache to prevent
// sending unnecessary requests to SofaScore.
func CollectMatches(db Database, team uint64) ([]Match, error) {
	ids, err := getEventIDs(team)
	if err != nil {
		return nil, err
	}

	var matches []Match
	for _, id := range ids {
		m, ok := db.FindMatch(id)
		if ok {
			matches = append(matches, m)
			continue
		}

		ev, err := getEvent(id)
		if err != nil {
			return nil, err
		}
		m = ev.toMatch()

		matches = append(matches, m)

		err = db.SaveMatch(m)
		if err != nil {
			fmt.Printf("Failed to save match to database: %s", err.Error())
		}
	}
	return matches, nil
}