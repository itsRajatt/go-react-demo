package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading env file")
	}
	PORT := os.Getenv("PORT")

	todos := []Todo{}

	app.Get("/api/todos", func(c fiber.Ctx) error {
		return c.Status(201).JSON(todos)
	})

	app.Post("/api/todos", func(c fiber.Ctx) error {
		todo := &Todo{}

		if err := c.Bind().Body(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body can't be empty"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	//UPDATE TODO
	app.Put("/api/todos/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	app.Delete("/api/todos/:id", func(c fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(todos)
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})

	})

	log.Fatal(app.Listen(":" + PORT))
}
