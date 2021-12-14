package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

// helper module needed for tests
// return test store with required configuration
// and return function which will cleanup the tables
func TestDB(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
	t.Helper() // TODO: read about this method

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	// returning DB and teardown function which will delete some tables
	return db, func(tables ...string) {
		// number of tables can be zero
		if len(tables) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		}

		db.Close()
	}
}
