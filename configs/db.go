package configs

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)


func InitDB() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", "./cotacoes.db")
    if err != nil {
        return nil, err
    }

    createTableQuery := `CREATE TABLE IF NOT EXISTS cotacoes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        bid TEXT,
        timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

    _, err = db.Exec(createTableQuery)
    if err != nil {
        return nil, err
    }

    return db, nil
}