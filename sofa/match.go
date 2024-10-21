package sofa

import (
	"fmt"
	"time"
)

// Match is a struct containing information about a match.
type Match struct {
	// ID is the ID of the event in SofaScore.
	ID int64 `json:"id"`
	// Tournament is the name of the tournament this match finds place in.
	Tournament string `json:"tournament"`
	// HomeTeam is the name of the home team.
	HomeTeam string `json:"home_team"`
	// AwayTeam is the name of the away team.
	AwayTeam string `json:"away_team"`
	// StartTime is the time the match will start.
	StartTime time.Time `json:"start_time"`
	// Location is the location the match will find place in.
	Location string `json:"location"`
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
			return nil, err
		}
	}
	return matches, nil
}
