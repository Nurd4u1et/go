package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	ID      int     `db:"id"`
	Name    string  `db:"name"`
	Email   string  `db:"email"`
	Balance float64 `db:"balance"`
}

func main() {
	db, err := sqlx.Open("postgres", "host=localhost port=5432 user=postgres password=1234 dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalln("Connection error:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalln("Failed to connect to the database:", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	fmt.Println("Connection to the database successful!")

	newUser := User{Name: "Samat", Email: "Samat@kbtu.com", Balance: 800.00}
	err = InsertUser(db, newUser)
	if err != nil {
		log.Println("Past error", err)
	}

	users, err := GetAllUsers(db)
	if err != nil {
		log.Println("Receiving error:", err)
	}

	fmt.Println("\nUser list:")
	for _, u := range users {
		fmt.Printf("%d | %s | %s | %.2f\n", u.ID, u.Name, u.Email, u.Balance)
	}
	err = TransferBalance(db, 1, 2, 100)
	if err != nil {
		log.Println("transfer error", err)
	} else {
		fmt.Println("\n transfer succes")
	}
}

func InsertUser(db *sqlx.DB, user User) error {
	query := `INSERT INTO users (name, email, balance) VALUES (:name, :email, :balance)`
	_, err := db.NamedExec(query, user)
	return err
}

func GetAllUsers(db *sqlx.DB) ([]User, error) {
	var users []User
	err := db.Select(&users, "SELECT * FROM users")
	return users, err
}

func GetUserByID(db *sqlx.DB, id int) (User, error) {
	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

func TransferBalance(db *sqlx.DB, fromID int, toID int, amount float64) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	var fromBalance float64
	err = tx.Get(&fromBalance, "SELECT balance FROM users WHERE id=$1", fromID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("User not found")
	}

	if fromBalance < amount {
		tx.Rollback()
		return fmt.Errorf("not enough funds")
	}

	_, err = tx.Exec("UPDATE users SET balance = balance - $1 WHERE id=$2", amount, fromID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE users SET balance = balance + $1 WHERE id=$2", amount, toID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
