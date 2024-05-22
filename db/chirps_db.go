package db

type Chirp struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	AuthID int    `json:"author_id"`
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
func (db *DB) GetChirps() ([]Chirp, error) {

	var chirps []Chirp

	// 执行查询
	rows, err := db.DataBase.Query("SELECT id, body, author_id FROM chirps")
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
