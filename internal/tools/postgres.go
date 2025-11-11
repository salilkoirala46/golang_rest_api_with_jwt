package tools

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
)

type postgresDB struct {
	db *sql.DB
}

var jwtkey = []byte("secret_key")

func (db *postgresDB) RegisterUser(user *UserDetail, password string) (*Token, error) {
	// ✅ Hash password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password", err)
		return nil, err
	}

	// Store the hashed password (not plain text)
	_, err = db.db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)",
		user.Name, user.Email, string(hashedPassword))
	if err != nil {
		log.Error("Registering user in PostgreSQL database failed", err)
		return nil, err
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		log.Error("Failed to sign JWT token", err)
		return nil, err
	}

	return &Token{Token: tokenString}, nil
}

func (db *postgresDB) Login(user *UserDetail, password string) (*Token, error) {
	// Implement login logic using PostgreSQL
	var storedHashedPassword string
	var name, email string

	err := db.db.QueryRow("SELECT name, email, password FROM users WHERE name=$1", user.Name).
		Scan(&name, &email, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error("Invalid username or password")
			return nil, nil
		}
		log.Error("Login query failed", err)
		return nil, err
	}

	// ✅ Compare the stored hash with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password)); err != nil {
		log.Error("Invalid password")
		return nil, nil
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		log.Error("Failed to sign JWT token", err)
		return nil, err
	}

	return &Token{Token: tokenString}, nil
}

func (db *postgresDB) ValidateToken(tokenString string) (*Token, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtkey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return &Token{Token: tokenString}, nil
}

func (db *postgresDB) SetupDatabase() error {
	// Implement DB setup/connection logic
	connStr := "postgresql://postgres:secret@localhost:5432/postgres?sslmode=disable"
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("Failed to establish connection ", err)
		return err
	}
	if err = database.Ping(); err != nil {
		log.Error("Error connecting to the database:", err)
		return err
	}
	db.db = database
	return nil
}

func (db *postgresDB) GetAllUsers() (*[]UserDetail, error) {
	// Query all users from PostgreSQL
	rows, err := db.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Error("Fetched users from PostgreSQL database failed", err)
		return nil, err
	}
	defer rows.Close()

	var users []UserDetail
	for rows.Next() {
		var user UserDetail
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}

func (db *postgresDB) AddUser(user *UserDetail) (*[]UserDetail, error) {
	// Insert user into PostgreSQL and return updated list
	_, err := db.db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		log.Error("Inserting user into PostgreSQL database failed", err)
		return nil, err
	}
	return db.GetAllUsers()
}

func (db *postgresDB) UpdateUser(updatedUser UserDetail, id int) (*[]UserDetail, error) {
	// Update user in PostgreSQL and return updated list
	_, err := db.db.Exec("UPDATE users SET name=$1, email=$2 WHERE id=$3", updatedUser.Name, updatedUser.Email, id)
	if err != nil {
		log.Error("Updating user in PostgreSQL database failed", err)
		return nil, err
	}
	return db.GetAllUsers()
}

func (db *postgresDB) DeleteUser(id int) (*[]UserDetail, error) {
	// Delete user from PostgreSQL and return updated list
	_, err := db.db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		log.Error("Deleting user from PostgreSQL database failed", err)
		return nil, err
	}
	return db.GetAllUsers()
}
