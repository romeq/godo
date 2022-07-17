package database

import (
	"database/sql"
	"errors"
	"os"
	"path"

    "github.com/romeq/godo/utils"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() error {
    // Get database path
    dbpath, err := getDbPath()
    if err != nil { return err }

    // Open database
    db, err = sql.Open("sqlite3", dbpath)
    if err != nil { return err }

    // Create connection
    err = db.Ping()
    if err != nil { return err }
    
    // Setup database
    _, err = setup(db)
    if err != nil { return err }
    
    return nil
}

func setup(db *sql.DB) (sql.Result, error) {
    return db.Exec(`CREATE TABLE IF NOT EXISTS todos(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        task VARCHAR(128),
        done BOOLEAN,
        deadline VARCHAR(128),
        date VARCHAR(128)
    )`)
}

func getDbPath() (string, error) {
    datadir := os.Getenv("XDG_DATA_HOME")
    if datadir == "" {
        homepath := os.Getenv("HOME")
        if homepath == "" {
            return "", errors.New("HOME environment variable doesn't exist")
        }
        datadir = path.Join(homepath, ".local/share")
    }

    projectdatadir := path.Join(datadir, "godo")
    utils.Check(os.MkdirAll(projectdatadir, 0755))
    return path.Join(projectdatadir, "godo.db"), nil
}

func ensureDb() {
    if db == nil {
        panic("Database not initialized before running queries!")
    }
}

