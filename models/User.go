package models

import (
	"crypto/rand"
	"fmt"
	"wastecontrol/db"

	"golang.org/x/crypto/scrypt"
)

type User struct {
	Id        int     `json:"id" binding:""`
	Username  string  `json:"username,omitempty" binding:""`
	Password  string  `json:"password,omitempty" binding:""`
	Token     string  `json:"token" binding:""`
	FirstName string  `json:"first_name" binding:""`
	LastName  string  `json:"last_name" binding:""`
	UserType  int     `json:"user_type" binding:""`
	Gdpr      int     `json:"gdpr" binding:""`
	Email     string  `json:"email" binding:""`
	Phone     *string `json:"phone" binding:""`
	Address   *string `json:"address" binding:""`
	Zip       *string `json:"zip" binding:""`
	City      *string `json:"city" binding:""`
}

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

/********** Void Functions **********/

func (u *User) CheckUsernameAndPassword() {
	checkLogin, _ := db.Conn.Query("SELECT id, user_type, gdpr FROM user WHERE username = ? AND password = ?", u.Username, EncryptPassword(u.Password))

	for checkLogin.Next() {
		checkLogin.Scan(&u.Id, &u.UserType, &u.Gdpr)

		u.GenerateToken()
	}
	defer checkLogin.Close()
}

func (u *User) GenerateToken() {
	u.Token = TokenGenerator()
	stmt, _ := db.Conn.Prepare("INSERT INTO auth (expiration, token, user_id) VALUES (NOW() + INTERVAL 1 DAY, ?, ?)")
	stmt.Exec(&u.Token, &u.Id)

	defer stmt.Close()
}

func (u *User) GetId() {
	getId, _ := db.Conn.Query("SELECT user_id FROM auth WHERE token = ?", u.Token)
	defer getId.Close()

	for getId.Next() {
		getId.Scan(&u.Id)
	}
}

func (c *User) Create(u *User) {
	stmt, _ := db.Conn.Prepare("INSERT INTO user (username, first_name, last_name, email, phone, address, postal, city, user_type, password, active, created_on) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 1, NOW())")
	res, _ := stmt.Exec(&u.Username, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.Address, &u.Zip, &u.City, &u.UserType, EncryptPassword(u.Password))

	id, _ := res.LastInsertId()
	stmt, _ = db.Conn.Prepare("INSERT INTO user_user (user_id, created_by) VALUES (?, ?)")
	stmt.Exec(&id, &c.Id)

	defer stmt.Close()
}

func (u *User) Get() {
	get, _ := db.Conn.Query("SELECT first_name, last_name, user_type FROM user WHERE id = ?", u.Id)

	for get.Next() {
		get.Scan(&u.FirstName, &u.LastName, &u.UserType)
	}
	defer get.Close()
}

func (u *User) GetRole() {
	get, _ := db.Conn.Query("SELECT user_type FROM user WHERE id = ?", u.Id)
	defer get.Close()

	for get.Next() {
		get.Scan(&u.UserType)
	}
}

func (u *User) GetSpecific() {
	db.Conn.QueryRow("SELECT username, first_name, last_name, email, phone, address, postal, city, user_type FROM user WHERE id = ?", u.Id).Scan(&u.Username, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.Address, &u.Zip, &u.City, &u.UserType)
}

func (u *User) Update() {
	stmt, _ := db.Conn.Prepare("UPDATE user SET username = ?,first_name = ?, last_name = ?, email = ?, phone = ?, address = ?, postal = ?, city = ?, user_type = ? WHERE id = ?")
	stmt.Exec(&u.Username, &u.FirstName, &u.LastName, &u.Email, &u.Phone, &u.Address, &u.Zip, &u.City, &u.UserType, &u.Id)

	defer stmt.Close()
}

func (u *User) UpdatePass() {
	stmt, _ := db.Conn.Prepare("UPDATE user SET password = ? WHERE id = ?")
	stmt.Exec(EncryptPassword(u.Password), &u.Id)

	defer stmt.Close()
}

func (u *User) UpdateGdpr() {
	stmt, _ := db.Conn.Prepare("UPDATE user SET gdpr = 1 WHERE id = ?")
	stmt.Exec(&u.Id)

	defer stmt.Close()
}

func (u *User) Delete() {
	stmt, _ := db.Conn.Prepare("UPDATE user SET deleted = 1 WHERE id = ?")
	stmt.Exec(&u.Id)

	defer stmt.Close()
}

/********** Return type Functions **********/

func EncryptPassword(password string) string {
	salt := make([]byte, PW_SALT_BYTES)
	hash, _ := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, PW_HASH_BYTES)

	return fmt.Sprintf("%x", hash)
}

func TokenGenerator() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func GetUsersForAdmin() []User {
	var users []User
	rows, _ := db.Conn.Query("SELECT id, username, first_name, last_name, user_type FROM user WHERE deleted != 1 ORDER BY id")
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.UserType)
		users = append(users, user)
	}
	defer rows.Close()

	return users
}

func GetUsersCreatedByUser(u *User) []User {
	var users []User
	rows, _ := db.Conn.Query("SELECT u.id, u.username, u.first_name, u.last_name, u.user_type FROM user u WHERE u.deleted != 1 AND u.id IN (SELECT user_id FROM user_user WHERE created_by = ?) ORDER BY u.id", u.Id)
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.UserType)
		users = append(users, user)
	}
	defer rows.Close()

	return users
}

func GetResellers() []User {
	var users []User
	rows, _ := db.Conn.Query("SELECT id, username FROM user WHERE user_type = 2")
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username)
		users = append(users, user)
	}
	defer rows.Close()

	return users
}

func GetOperators() []User {
	var users []User
	rows, _ := db.Conn.Query("SELECT id, username FROM user WHERE user_type = 3")
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username)
		users = append(users, user)
	}
	defer rows.Close()

	return users
}

func GetViewers(role int, userId int) []User {
	var users []User
	rows, _ := db.Conn.Query("SELECT id, username FROM user WHERE deleted = 0 AND user_type = 4")
	defer rows.Close()
	if role != 1 {
		rows, _ = db.Conn.Query("SELECT u.id, u.username FROM user u LEFT JOIN user_user uu ON uu.user_id = u.id WHERE u.deleted = 0 AND uu.created_by = ? AND u.user_type = 4", userId)
		defer rows.Close()
	}
	for rows.Next() {
		var user User
		rows.Scan(&user.Id, &user.Username)
		users = append(users, user)
	}

	return users
}
