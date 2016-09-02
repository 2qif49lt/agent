package eventdb

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// 所有请求的操作agentd都应生存一个guid mission id 和一个本地自增id

const (
	// mission 事件
	createMissionsTable = `
    CREATE TABLE IF NOT EXISTS missions (
    	id integer PRIMARY KEY autoincrement,
        mid text NOT NULL,
        version text NOT NULL,
        command text NOT NULL,
        paras text NOT NULL,
        body text,
        result text NOT NULL DEFAULT 'OK',
        cost integer NOT NULL,
        logtime DEFAULT (datetime('now','localtime')),
        CONSTRAINT key_unique UNIQUE (mid)
    );
    `
	// 通用事件
	createLogsTable = `
    `

	// plugin 事件
	createPluginsTable = `
    `
)

type EventRec struct {
	Id      int       `json:"id"`
	Mid     string    `json:"mid"`
	Version string    `json:"version"`
	Command string    `json:"command"`
	Paras   string    `json:"paras"`
	Body    string    `json:"body"`
	Result  string    `json:"result"`
	Cost    int       `json:"cost"`
	Logtime time.Time `json:"logtime"`
}

// Database is a graph database for storing entities and their relationships.
type Database struct {
	conn *sql.DB
	mux  sync.RWMutex
}

// NewDatabase creates a new graph database initialized with a root entity.
func NewDatabase(conn *sql.DB) (*Database, error) {
	if conn == nil {
		return nil, fmt.Errorf("Database connection cannot be nil")
	}
	db := &Database{conn: conn}

	if _, err := conn.Exec(createMissionsTable); err != nil {
		return nil, err
	}

	return db, nil
}

func Close() error {
	err := eventer.Close()
	return err
}

// Close the underlying connection to the database.
func (db *Database) Close() error {
	return db.conn.Close()
}

func InsertMission(mid, version, command, paras, body, result string, cost int) error {
	return eventer.InsertMission(mid, version, command, paras, body, result, cost)
}
func (db *Database) InsertMission(mid, version, command, paras, body, result string, cost int) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	_, err := db.conn.Exec(`INSERT INTO missions(mid,version,command,paras,body,result,cost) VALUES(?,?,?,?,?,?,?);`,
		mid, version, command, paras, body, result, cost)

	return err
}
func GetMission(mid string) (*EventRec, error) {
	return eventer.GetMission(mid)
}
func (db *Database) GetMission(mid string) (*EventRec, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	rec := &EventRec{}
	rec.Mid = mid

	if err := db.conn.QueryRow(`SELECT id,version,command,paras,body,result,cost,logtime FROM missions WHERE mid = ?;`, mid).Scan(&rec.Id,
		rec.Version, rec.Command, rec.Paras, rec.Body, rec.Result, rec.Cost, rec.Logtime); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return rec, nil
}
