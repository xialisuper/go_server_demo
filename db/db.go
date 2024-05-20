package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"  // 导入 pq 包

)

type DB struct {
	path     string
	DataBase *sql.DB
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
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

// CreateChirp creates a new chirp and saves it to database
func (db *DB) CreateChirp(body string) (Chirp, error) {

	// 插入chirp到数据库
	var chirp Chirp
	err := db.DataBase.QueryRow(
		"INSERT INTO chirps (body) VALUES ($1) RETURNING id, body",
		body,
	).Scan(&chirp.ID, &chirp.Body)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

	var chirps []Chirp

	// 执行查询
	rows, err := db.DataBase.Query("SELECT id, body FROM chirps")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var chirp Chirp
		err = rows.Scan(&chirp.ID, &chirp.Body)
		if err != nil {
			return nil, err
		}
		chirps = append(chirps, chirp)
	}

	// 检查是否有查询错误
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return chirps, nil

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
