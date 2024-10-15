package handler

import (
	"fmt"
	"openruntimes/handler/sofa"
	"os"
	"strconv"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/open-runtimes/types-for-go/v4/openruntimes"
)

const DB_NAME = "sofadata"

func Main(Context openruntimes.Context) openruntimes.Response {
	client := appwrite.NewClient(
		appwrite.WithEndpoint(os.Getenv("APPWRITE_API_ENDPOINT")),
		appwrite.WithProject(os.Getenv("APPWRITE_PROJECT_ID")),
		appwrite.WithKey(os.Getenv("APPWRITE_SECRET_API_KEY")),
	)
	databases := appwrite.NewDatabases(client)

	// Create new DB.
	db := NewDB(databases)

	// Fetch all teams.
	docs, err := databases.ListDocuments(os.Getenv("APPWRITE_DB_ID"), "teams",
		databases.WithListDocumentsQueries([]string{
			"{\"method\":\"select\",\"values\":[]}",
		}))
	if err != nil {
		Context.Error(err)
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
			Context.Log(err.Error())
			continue
		}
		Context.Log(fmt.Sprintf("Collected %d matches for team %s\n", len(matches), d.Id))
		sum += len(matches)
	}

	return Context.Res.Text(fmt.Sprintf("Collected or refreshed %d matches.", sum))
}
