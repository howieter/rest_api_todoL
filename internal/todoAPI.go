package todoAPI

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func SetupRoutes(app *fiber.App, db *pgx.Conn) {
	app.Get("/tasks", func(c *fiber.Ctx) error {
		rows, err := db.Query(context.Background(), "SELECT id, title, description, status FROM tasks")
		if err != nil {
			return c.Status(500).SendString("db query error")
		}
		defer rows.Close()

		var tasks []map[string]interface{}
		for rows.Next() {
			var id int
			var title, description, status string
			if err := rows.Scan(&id, &title, &description, &status); err != nil {
				return c.Status(500).SendString("data processing error")
			}
			tasks = append(tasks, map[string]interface{}{
				"id":          id,
				"title":       title,
				"description": description,
				"status":      status,
			})
		}

		return c.JSON(fiber.Map{"tasks": tasks})
	})

	app.Post("/tasks", func(c *fiber.Ctx) error {
		var task struct {
			Title       string `json:"title"`
			Description string `json:"description"`
		}

		if err := c.BodyParser(&task); err != nil {
			return c.Status(400).SendString("error parsing JSON")
		}

		_, err := db.Exec(context.Background(),
			"INSERT INTO tasks (title, description, status) VALUES ($1, $2, 'new')",
			task.Title, task.Description)

		if err != nil {
			return c.Status(500).SendString("error insert to bd")
		}

		return c.JSON(fiber.Map{"message": "task created"})
	})

	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var task struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Status      string `json:"status"`
		}

		if err := c.BodyParser(&task); err != nil {
			return c.Status(400).SendString("error parsing JSON")
		}

		_, err := db.Exec(context.Background(),
			"UPDATE tasks SET title=$1, description=$2, status=$3 WHERE id=$4",
			task.Title, task.Description, task.Status, id)

		if err != nil {
			return c.Status(500).SendString("error update task")
		}

		return c.JSON(fiber.Map{"message": "task " + id + " updated"})
	})

	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		_, err := db.Exec(context.Background(), "DELETE FROM tasks WHERE id=$1", id)
		if err != nil {
			return c.Status(500).SendString("error task delete")
		}

		return c.JSON(fiber.Map{"message": "task " + id + " deleted"})
	})
}
