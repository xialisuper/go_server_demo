package db

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Expires  int64  `json:"expires_in_seconds"`
}

// LoginUser 登录用户
func (db *DB) LoginUser(email string, password string) (User, error) {
	var user User
	err := db.DataBase.QueryRow(
		"SELECT id, email, password FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser 创建一个新用户并保存至数据库
func (db *DB) CreateUser(email string, password string) (User, error) {
	var user User

	hashedPassword, err := GenerateFromPassword(password)

	if err != nil {
		return User{}, err
	}

	err = db.DataBase.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, password",
		email,
		hashedPassword,
	).Scan(&user.ID, &user.Email, &user.Password)

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

func (db *DB) UpdateUser(id int, email string, password string) (User, error) {
	var user User
	hashedPassword, err := GenerateFromPassword(password)
	if err != nil {
		return User{}, err
	}
	err = db.DataBase.QueryRow(
		"UPDATE users SET email = $1, password = $2 WHERE id = $3 RETURNING id, email",
		email,
		hashedPassword,
		id,
	).Scan(&user.ID, &user.Email)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GenerateFromPassword  hash password
func GenerateFromPassword(password string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// CompareHashAndPassword  compare hashed password
func CompareHashAndPassword(hashedPassword, password string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

}
