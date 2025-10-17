package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
)

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
}