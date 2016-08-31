package eventdb

import (
	"database/sql"
	"path"
)

const (
	dbfolder    = "db"
	eventdbname = "enent.db"
)

var (
	eventer *Database = nil
)

func init() {
	db, err := NewSqliteConn(path.Join(dbfolder, eventdbname))
	if err != nil {
		panic(err)
	}
	eventer = db
}

// NewSqliteConn opens a connection to a sqlite
// database.
func NewSqliteConn(root string) (*Database, error) {
	conn, err := sql.Open("sqlite3", root)
	if err != nil {
		return nil, err
	}
	return NewDatabase(conn)
}
