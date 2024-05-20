package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // 导入 pq 包
)

type DB struct {
	path     string
	DataBase *sql.DB
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {

	// 连接到数据库
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, err
	}

	myDb := DB{
		path:     path,
		DataBase: db,
	}

	err = myDb.ensureDB()

	if err != nil {
		return nil, err
	}

	return &myDb, nil
}


// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {

	// 验证连接
	err := db.DataBase.Ping()
	if err != nil {
		return err
	}
	fmt.Println("Successfully connected to the database!")

	return nil
}
