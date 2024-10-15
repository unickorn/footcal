package handler

import (
	"fmt"
	"openruntimes/handler/sofa"
	"os"

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
	type team struct {
		ID uint64 `json:"id"`
	}
	docs, err := databases.ListDocuments(os.Getenv("APPWRITE_DB_ID"), "teams", 
		databases.WithListDocumentsQueries([]string {
			"{\"method\":\"select\",\"values\":[\"id\"]}",
		}))
	if err != nil {
		Context.Error(err)
		return Context.Res.Empty()
	}

	for _, d := range docs.Documents {
		var t team

		err = d.Decode(&t)
		if err != nil {
			Context.Error(err)
			return Context.Res.Empty()
		}

		// Fetch all events of the team, saving ones that cannot be found.
		matches, err := sofa.CollectMatches(db, t.ID)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Printf("Collected %d matches for team %d\n", len(matches), t.ID)
	}

	return Context.Res.Text("Pong")
}
