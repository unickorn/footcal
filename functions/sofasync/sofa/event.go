package sofa

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type event struct {
	ID         int64  `json:"id"`
	Slug       string `json:"slug"`
	Tournament struct {
		Name string `json:"name"`
	} `json:"tournament"`
	HomeTeam struct {
		Name string `json:"name"`
	} `json:"homeTeam"`
	AwayTeam struct {
		Name string `json:"name"`
	} `json:"awayTeam"`
	StartTime uint64 `json:"startTimestamp"`
	Venue     struct {
		City struct {
			Name string `json:"name"`
		} `json:"city"`
		Stadium struct {
			Name string `json:"name"`
		} `json:"stadium"`
	} `json:"venue"`
}

func (ev event) toMatch() Match {
	return Match{
		Tournament: ev.Tournament.Name,
		HomeTeam:   ev.HomeTeam.Name,
		AwayTeam:   ev.AwayTeam.Name,
		StartTime:  time.Unix(int64(ev.StartTime), 0),
		Location:   fmt.Sprintf("%s, %s", ev.Venue.City.Name, ev.Venue.Stadium.Name),
	}
}

func getEventIDs(team uint64) ([]int64, error) {
	rsp, err := http.Get(fmt.Sprintf(eventsEndpoint, team))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	// Parse the response into a struct.
	var eventIDs struct {
		Events []struct {
			ID int64 `json:"id"`
		} `json:"events"`
	}
	err = json.NewDecoder(rsp.Body).Decode(&eventIDs)
	if err != nil {
		return nil, err
	}
	// Return the IDs.
	var ids []int64
	for _, e := range eventIDs.Events {
		ids = append(ids, e.ID)
	}
	return ids, nil
}

func getEvent(id int64) (event, error) {
	rsp, err := http.Get(fmt.Sprintf(eventEndpoint, id))
	if err != nil {
		return event{}, err
	}
	defer rsp.Body.Close()

	var ev event
	err = json.NewDecoder(rsp.Body).Decode(&ev)
	if err != nil {
		return event{}, err
	}
	return ev, nil
}
