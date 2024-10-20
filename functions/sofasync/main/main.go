package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/appwrite/sdk-for-go/appwrite"
	"github.com/unickorn/footcal/sofa"
)

func main() {
	client := appwrite.NewClient(
		appwrite.WithEndpoint(os.Getenv("APPWRITE_API_ENDPOINT")),
		appwrite.WithProject(os.Getenv("APPWRITE_PROJECT_ID")),
		appwrite.WithKey(os.Getenv("APPWRITE_SECRET_API_KEY")),
	)
	databases := appwrite.NewDatabases(client)

	// Create new DB.
	db := DB{}

	// Fetch all teams.
	docs, err := databases.ListDocuments(os.Getenv("APPWRITE_DB_ID"), "teams")
	if err != nil {
		fmt.Println(err)
		return
	}

	var sum int
	for _, d := range docs.Documents {
		id, err := strconv.ParseUint(d.Id, 10, 64)
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to parse team ID %s", d.Id))
			return
		}
		// Fetch all events of the team, saving ones that cannot be found.
		matches, err := sofa.CollectMatches(db, id)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println(matches)
		fmt.Println(fmt.Sprintf("Collected %d matches for team %s\n", len(matches), d.Id))
		sum += len(matches)
	}

	fmt.Println(fmt.Sprintf("Collected or refreshed %d matches.", sum))
}
