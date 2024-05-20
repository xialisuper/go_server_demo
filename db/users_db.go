package db

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// CreateUser 创建一个新用户并保存至数据库
func (db *DB) CreateUser(email string) (User, error) {
	var user User
	err := db.DataBase.QueryRow(
		"INSERT INTO users (email) VALUES ($1) RETURNING id, email",
		email,
	).Scan(&user.ID, &user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// GetUserByID 根据 id 返回一个用户
func (db *DB) GetUserByID(id int) (User, error) {
	var user User
	err := db.DataBase.QueryRow(
		"SELECT id, email FROM users WHERE id = $1",
		id,
	).Scan(&user.ID, &user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// GetUsers 返回数据库中的所有用户
func (db *DB) GetUsers() ([]User, error) {
	var users []User
	rows, err := db.DataBase.Query("SELECT id, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
