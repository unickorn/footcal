package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/appwrite/sdk-for-go/appwrite"
	ics "github.com/arran4/golang-ical"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
	"github.com/unickorn/footcal/sofa"
)

func Main(Context openruntimes.Context) openruntimes.Response {
	client := appwrite.NewClient(
		appwrite.WithProject(os.Getenv("APPWRITE_PROJECT_ID")),
		appwrite.WithKey(Context.Req.Headers["x-appwrite-key"]),
	)
	databases := appwrite.NewDatabases(client)

	if Context.Req.Method == "GET" {
		Context.Log(fmt.Sprintf("%#+v", Context.Req.Query))
		teams, ok := Context.Req.Query["teamlist"]
		if !ok {
			Context.Error("No teams specified!")
			return Context.Res.Text("No teams specified!",
				Context.Res.WithStatusCode(http.StatusBadRequest))
		}

		// teams := "3052"
		// Parse teams.
		teamsSplitted := strings.Split(teams, ",")

		// Collect matches.
		var matches []sofa.Match

		for _, t := range teamsSplitted {
			_, err := strconv.ParseUint(t, 10, 64)
			if err != nil {
				Context.Error(err.Error())
				return Context.Res.Text("Error parsing uint: " + err.Error())
			}

			// Get team names from DB:
			Context.Log("Getting db " + os.Getenv("APPWRITE_DB_ID") + " table teams with ID " + t)

			docs, err := databases.ListDocuments(os.Getenv("APPWRITE_DB_ID"), "teams")
			if err != nil 
				Context.Error(err.Error())
				return Context.Res.Text("Error getting team names for team " + t + ": " + err.Error())
			}
			for _, d := range docs.Documents {
				Context.Log(fmt.Sprint("Doc found with ID: %#+v", d.Id))
			}
			doc, err := databases.GetDocument(os.Getenv("APPWRITE_DB_ID"), "teams", t)
			if err != nil {
				Context.Error(err.Error())
				return Context.Res.Text("Error getting team names for team " + t + ": " + err.Error())
			}
			type team struct {
				Name string `json:"name"`
			}
			var t team
			err = doc.Decode(&t)
			if err != nil {
				Context.Error(err.Error())
				return Context.Res.Text("Error decoding teams: " + err.Error())
			}

			// Get matches from team name now.
			name := t.Name

			list, err := databases.ListDocuments(os.Getenv("APPWRITE_DB_ID"), "matches", databases.WithListDocumentsQueries(
				[]string{
					Or([]string{Equal("home_team", name), Equal("away_team", name)}),
				},
			))
			if err != nil {
				Context.Error(err.Error())
				return Context.Res.Text("Error getting matches from team name:" + err.Error())
			}

			for _, m := range list.Documents {
				var match sofa.Match
				err = m.Decode(&match)
				if err != nil {
					Context.Error(err.Error())
					return Context.Res.Text("Error decoding listed matches: " + err.Error())
				}

				matches = append(matches, match)
			}
		}

		// Convert all collected matches to event and add to calendar.
		cal := ics.NewCalendar()
		cal.SetMethod(ics.MethodPublish)

		for _, match := range matches {
			event := cal.AddEvent(fmt.Sprintf("%s - %s", match.HomeTeam, match.AwayTeam))
			event.SetDtStampTime(time.Now())
			event.SetStartAt(match.StartTime)
			event.SetEndAt(match.StartTime.Add(time.Hour * 2))
			event.SetSummary(fmt.Sprintf("%s - %s", match.HomeTeam, match.AwayTeam))
			event.SetDescription(match.Tournament)
			event.SetLocation(match.Location)
		}
		data := cal.Serialize()
		// Return the calendar.
		return Context.Res.Binary([]byte(data), Context.Res.WithHeaders(map[string]string{
			"Content-Type":        "text/calendar",
			"Content-Disposition": "attachment; filename=calendar.ics",
		}))
	}

	return Context.Res.Text("Not a GET request")
}
