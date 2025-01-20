package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DataSourceName string = "./synapse.db"
	SQLDir         string = "data/sql"
)

const (
	MigrationApplied  MigrationState = iota
	MigrationPending  MigrationState = iota
	MigrationFailed   MigrationState = iota
	MigrationReverted MigrationState = iota
)

// MigrationState represents an enumerable description of the
// state of a migration. See [MigrationState.String] for mapping.
type MigrationState int

type Repo interface {
	GetDB() (conn *sql.DB)
	InsertRow(rec Repo) (ok bool, id string, err error)
}

type Connection struct {
	Db *sql.DB
}

// struct Migration represents a single [sql.DB] migration
type Migration struct {
	ID      string
	SQL     string
	Applied string
}

type Schema struct {
	Version     int
	State       int
	Description string
	AppliedAt   time.Time
	Hash        string
}

type Metadata struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

func (s MigrationState) String() string {
	switch s {
	case MigrationApplied:
		return "applied"
	case MigrationFailed:
		return "failed"
	case MigrationReverted:
		return "rolled_back"
	default:
		return "pending"
	}
}

func CreateConnection() *Connection {
	var err error
	c := Connection{}

	if c.Db, err = sql.Open("sqlite3", DataSourceName); err != nil {
		logger.Fatal(
			fmt.Sprintf("unable to connect to database %v %v", DataSourceName, err.Error()),
		)
	}

	return &c
}

func (c Connection) ExecuteSQL(fpath string) error {
	sb := strings.Builder{}
	contents, err := os.ReadFile(fpath)
	if err != nil {
		return fmt.Errorf("unable to read file %v %v", fpath, err.Error())
	}

	sb.Write(contents)

	fileContents := sb.String()
	logger.Debugf("query:\n%v\n...", strings.Join(strings.Split(fileContents, "\n")[:3], "\n"))

	res, err := c.Db.Exec(fileContents)
	if err != nil {
		return fmt.Errorf("unable to execute sql %v %v", fpath, err.Error())
	}

	if id, err := res.LastInsertId(); err != nil {
		return err
	} else {
		logger.Info(fmt.Sprintf("ID of last inserted row %v", id))
		return nil
	}
}

// function Setup is a first-run or forced recreation of the database
//
// TODO: Replace with a prompt
func SetupDb(force bool) error {
	pending := []string{}
	if force {
		logger.Warn("you're about to delete the database! at " + DataSourceName)
		os.Remove(DataSourceName)
	}

	dirEntries, err := os.ReadDir(SQLDir)
	if err != nil {
		logger.Errorf("unable to read dir %v/ %v", SQLDir, err.Error())
	}

	for _, e := range dirEntries {
		inf, _ := e.Info()
		fname := inf.Name()

		logger.Debugf("file: %v (%v bytes)", fname, inf.Size())

		if strings.HasSuffix(fname, ".sql") {
			pending = append(pending, fmt.Sprintf("%v/%v", SQLDir, fname))
		}
	}

	c := CreateConnection()

	for _, f := range pending {
		err = c.ExecuteSQL(f)
		if err != nil {
			logger.Errorf("query failed with err %v", err.Error())

			return err
		}

		s, _ := strings.CutSuffix(f, ".sql")
		logger.Infof("created table %v", s)
	}

	return nil
}
