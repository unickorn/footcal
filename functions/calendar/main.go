package handler

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/arran4/golang-ical"
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
		teams, ok := Context.Req.Query["teams"]
		if !ok {
			Context.Error("No teams specified!")
			return Context.Res.Text("No teams specified!",
				Context.Res.WithStatusCode(http.StatusBadRequest))
		}

		// Parse teams.
		teamsSplitted := strings.Split(teams, ",")

		// Collect matches.
		var matches []sofa.Match

		for _, t := range teamsSplitted {
			_, err := strconv.ParseUint(t, 10, 64)
			if err != nil {
				Context.Error(err.Error())
			}

			// Get team names from DB:
			doc, err := databases.GetDocument(os.Getenv("APPWRITE_DB_NAME"), "teams", t)
			if err != nil {
				Context.Error(err.Error())
				return Context.Res.Empty()
			}
			type team struct {
				Name string `json:"name"`
			}
			var t team
			err = doc.Decode(&t)
			if err != nil {
				Context.Error(err.Error())
				return Context.Res.Empty()
			}

			// Get matches from team name now.
			name := t.Name

			list, err := databases.ListDocuments(os.Getenv("APPWRITE_DB_NAME"), "matches", databases.WithListDocumentsQueries(
				[]string{
					Or([]string{Equal("home_team", name), Equal("away_team", name)}),
				},
			))
			if err != nil {
				Context.Error(err.Error())
				return Context.Res.Empty()
			}

			for _, m := range list.Documents {
				var match sofa.Match
				err = m.Decode(&match)
				if err != nil {
					Context.Error(err.Error())
					return Context.Res.Empty()
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
		return Context.Res.Binary([]byte(data), Context.Res.WithHeaders(map[string]string {
			"Content-Type": "text/calendar",
			"Content-Disposition": "attachment; filename=calendar.ics",
		}))
	}

	return Context.Res.Empty()
}
