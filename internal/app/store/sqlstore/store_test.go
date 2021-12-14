package sqlstore_test

import (
	"fmt"
	"os"
	"testing"
)

var databaseURL string

// TestMain is called once before all the tests
func TestMain(m *testing.M) {
	fmt.Println("Starting the tests...")
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=restapi_test user=postgres password=qwe123QWE sslmode=disable"
	}

	// TODO: read the docs for string below
	os.Exit(m.Run()) // Need to exit with correct code
}
