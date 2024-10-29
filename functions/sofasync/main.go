package handler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
	"github.com/unickorn/footcal/sofa"
)

func Main(Context openruntimes.Context) openruntimes.Response {
	client := appwrite.NewClient(
		appwrite.WithProject(os.Getenv("APPWRITE_PROJECT_ID")),
		appwrite.WithKey(Context.Req.Headers["x-appwrite-key"]),
	)
	databases := appwrite.NewDatabases(client)

	// Create new DB.
	db := NewDB(databases, Context)

	// Fetch all teams.
	docs, err := databases.ListDocuments(os.Getenv("APPWRITE_DB_ID"), "teams")
	if err != nil {
		Context.Error("Error fetching all teams: " + err.Error())
		return Context.Res.Empty()
	}

	var sum int
	for _, d := range docs.Documents {
		id, err := strconv.ParseUint(d.Id, 10, 64)
		if err != nil {
			Context.Error(fmt.Sprintf("Failed to parse team ID %s", d.Id))
			return Context.Res.Text("Error")
		}
		// Fetch all events of the team, saving ones that cannot be found.
		matches, err := sofa.CollectMatches(db, id)
		if err != nil {
			Context.Log("Failed to collect matches: " + err.Error())
			continue
		}
		Context.Log(fmt.Sprintf("Collected %d matches for team %s\n", len(matches), d.Id))
		sum += len(matches)
	}

	return Context.Res.Text(fmt.Sprintf("Collected or refreshed %d matches.", sum))
}
