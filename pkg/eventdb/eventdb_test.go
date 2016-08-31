package eventdb

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newTestDb(t *testing.T) (*Database, string) {
	p := path.Join(os.TempDir(), "sqlite.db")
	conn, err := sql.Open("sqlite3", p)
	db, err := NewDatabase(conn)
	if err != nil {
		t.Fatal(err)
	}
	return db, p
}

func destroyTestDb(dbPath string) {
	os.Remove(dbPath)
}

func TestNewDatabase(t *testing.T) {
	db, dbpath := newTestDb(t)
	if db == nil {
		t.Fatal("Database should not be nil")
	}
	db.Close()
	defer destroyTestDb(dbpath)
}
