package db

import "fmt"

type Chirp struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	AuthID int    `json:"author_id"`
}

// GetChirpsByAuthorID returns all chirps by author id
func (db *DB) GetChirpsByAuthorID(userID int, sort string) ([]Chirp, error) {
	if sort == "desc" {
		sort = "DESC"
	} else {
		sort = "ASC"
	}

	var chirps []Chirp

	// 执行查询 并以Continue sorting the chirps by id in "sort" order.
	rows, err := db.DataBase.Query(
		"SELECT id, body, author_id FROM chirps WHERE author_id = $1 ORDER BY id "+sort,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var chirp Chirp
		err = rows.Scan(&chirp.ID, &chirp.Body, &chirp.AuthID)
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

// DeleteChirpByID deletes a single chirp by id
func (db *DB) DeleteChirpByID(id int, userID int) error {
	// 执行删除
	result, err := db.DataBase.Exec("DELETE FROM chirps WHERE id = $1 AND author_id = $2", id, userID)
	if err != nil {
		return err
	}

	// 检查受影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no chirp found with id %d for user %d", id, userID)
	}

	return nil
}

// CreateChirp creates a new chirp and saves it to database
func (db *DB) CreateChirp(body string, userID int) (Chirp, error) {

	// 插入chirp到数据库
	var chirp Chirp
	err := db.DataBase.QueryRow(
		// "INSERT INTO chirps (body) VALUES ($1) RETURNING id, body",
		"INSERT INTO chirps (body, author_id) VALUES ($1, $2) RETURNING id, body, author_id",
		body, userID,
	).Scan(&chirp.ID, &chirp.Body, &chirp.AuthID)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil

}

// GetChirpByID returns a single chirp by id
func (db *DB) GetChirpByID(id int) (Chirp, error) {

	var chirp Chirp

	// 执行查询
	err := db.DataBase.QueryRow(
		"SELECT id, body, author_id FROM chirps WHERE id = $1",
		id,
	).Scan(&chirp.ID, &chirp.Body, &chirp.AuthID)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps(sort string) ([]Chirp, error) {

	if sort == "desc" {
		sort = "DESC"
	} else {
		sort = "ASC"
	}

	var chirps []Chirp

	// 执行查询 并以Continue sorting the chirps by id in "sort" order.
	rows, err := db.DataBase.Query(
		"SELECT id, body, author_id FROM chirps ORDER BY id "+sort,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 遍历查询结果
	for rows.Next() {
		var chirp Chirp
		err = rows.Scan(&chirp.ID, &chirp.Body, &chirp.AuthID)
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
