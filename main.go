package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "./mydb.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	createUsersTable := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		membership TEXT,
		name TEXT,
		surname TEXT,
		phone TEXT,
		email TEXT,
		join_date TEXT,
		membership_level TEXT,
		points INTEGER
	)`

	createTransfersTable := `CREATE TABLE IF NOT EXISTS transfers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_user_id INTEGER NOT NULL,
		to_user_id INTEGER NOT NULL,
		amount INTEGER NOT NULL CHECK (amount > 0),
		status TEXT NOT NULL CHECK (status IN ('pending','processing','completed','failed','cancelled','reversed')),
		note TEXT,
		idempotency_key TEXT NOT NULL UNIQUE,
		created_at TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		completed_at TEXT,
		fail_reason TEXT,
		FOREIGN KEY (from_user_id) REFERENCES users(id),
		FOREIGN KEY (to_user_id) REFERENCES users(id)
	)`

	createPointLedgerTable := `CREATE TABLE IF NOT EXISTS point_ledger (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		change INTEGER NOT NULL,
		balance_after INTEGER NOT NULL,
		event_type TEXT NOT NULL CHECK (event_type IN ('transfer_out','transfer_in','adjust','earn','redeem')),
		transfer_id INTEGER,
		reference TEXT,
		metadata TEXT,
		created_at TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (transfer_id) REFERENCES transfers(id)
	)`

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createTransfersTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createPointLedgerTable)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func setupRoutes(app *fiber.App, db *sql.DB) {
	app.Get("/users", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		users := []map[string]interface{}{}
		for rows.Next() {
			var id int
			var membership, name, surname, phone, email, join_date, membership_level string
			var points int
			rows.Scan(&id, &membership, &name, &surname, &phone, &email, &join_date, &membership_level, &points)

			users = append(users, map[string]interface{}{
				"id":              id,
				"membership":      membership,
				"name":            name,
				"surname":         surname,
				"phone":           phone,
				"email":           email,
				"join_date":       join_date,
				"membership_level": membership_level,
				"points":          points,
			})
		}
		return c.JSON(users)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		row := db.QueryRow("SELECT * FROM users WHERE id = ?", id)

		var membership, name, surname, phone, email, join_date, membership_level string
		var points int
		var userID int
		if err := row.Scan(&userID, &membership, &name, &surname, &phone, &email, &join_date, &membership_level, &points); err != nil {
			return c.Status(404).SendString("User not found")
		}

		user := map[string]interface{}{
			"id":              userID,
			"membership":      membership,
			"name":            name,
			"surname":         surname,
			"phone":           phone,
			"email":           email,
			"join_date":       join_date,
			"membership_level": membership_level,
			"points":          points,
		}
		return c.JSON(user)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		user := new(map[string]interface{})
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := "INSERT INTO users (membership, name, surname, phone, email, join_date, membership_level, points) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		result, err := db.Exec(query, (*user)["membership"], (*user)["name"], (*user)["surname"], (*user)["phone"], (*user)["email"], (*user)["join_date"], (*user)["membership_level"], (*user)["points"])
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		id, _ := result.LastInsertId()
		return c.JSON(map[string]interface{}{"id": id})
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		user := new(map[string]interface{})
		if err := c.BodyParser(user); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := "UPDATE users SET membership = ?, name = ?, surname = ?, phone = ?, email = ?, join_date = ?, membership_level = ?, points = ? WHERE id = ?"
		_, err := db.Exec(query, (*user)["membership"], (*user)["name"], (*user)["surname"], (*user)["phone"], (*user)["email"], (*user)["join_date"], (*user)["membership_level"], (*user)["points"], id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendStatus(200)
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendStatus(200)
	})

	app.Post("/point_ledger", func(c *fiber.Ctx) error {
		entry := new(map[string]interface{})
		if err := c.BodyParser(entry); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		query := "INSERT INTO point_ledger (user_id, change, balance_after, event_type, transfer_id, reference, metadata, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		result, err := db.Exec(query, (*entry)["user_id"], (*entry)["change"], (*entry)["balance_after"], (*entry)["event_type"], (*entry)["transfer_id"], (*entry)["reference"], (*entry)["metadata"], (*entry)["created_at"])
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		id, _ := result.LastInsertId()
		return c.JSON(map[string]interface{}{"id": id})
	})

	app.Get("/point_ledger", func(c *fiber.Ctx) error {
		userID := c.Query("user_id")
		rows, err := db.Query("SELECT * FROM point_ledger WHERE user_id = ?", userID)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()

		entries := []map[string]interface{}{}
		for rows.Next() {
			var id, user_id, change, balance_after, transfer_id int
			var event_type, reference, metadata, created_at string
			rows.Scan(&id, &user_id, &change, &balance_after, &event_type, &transfer_id, &reference, &metadata, &created_at)

			entries = append(entries, map[string]interface{}{
				"id":           id,
				"user_id":      user_id,
				"change":       change,
				"balance_after": balance_after,
				"event_type":   event_type,
				"transfer_id":  transfer_id,
				"reference":    reference,
				"metadata":     metadata,
				"created_at":   created_at,
			})
		}
		return c.JSON(entries)
	})
}

func main() {
	db := initDatabase()
	defer db.Close()

	app := fiber.New()
	setupRoutes(app, db)

	app.Listen(":3000")
}