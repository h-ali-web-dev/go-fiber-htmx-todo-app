package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("store.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	err = database.AutoMigrate(&Todo{})
	if err != nil {
		return
	}
	DB = database
	var todo []Todo
	DB.Find(&todo)
	if len(todo) < 1 {
		todos := []Todo{
			{Title: "todo 1", Completed: false, CreatedAt: time.Now(), ID: uuid.New()},
			{Title: "todo 2", Completed: true, CreatedAt: time.Now(), ID: uuid.New()},
		}
		DB.Create(&todos)
	}
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	ConnectDatabase()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Golang Htmx Todo List",
		})
	})

	app.Get("/todos", func(c *fiber.Ctx) error {
		var todo []Todo
		DB.Find(&todo)
		return c.Render("todo-list", fiber.Map{
			"Todos": todo,
		})
	})

	app.Post("/todos", func(c *fiber.Ctx) error {
		title := c.FormValue("Title")
		var completed bool = false
		if c.FormValue("Completed") == "on" {
			completed = true
		}
		todo := Todo{ID: uuid.New(), Title: title, Completed: completed, CreatedAt: time.Now()}
		DB.Create(&todo)
		var todoNew []Todo
		DB.Find(&todoNew)
		return c.Render("todo-list", fiber.Map{
			"Todos": todoNew,
		})
	})

	app.Put("/todos/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("ID"))
		if err != nil {
			return err
		}
		var user Todo
		DB.First(&user, id)
		userModel := DB.Model(&Todo{}).Where("id = ?", id)
		userModel.Update("Completed", !user.Completed)
		var todoNew []Todo
		DB.Find(&todoNew)
		return c.Render("todo-list", fiber.Map{
			"Todos": todoNew,
		})
	})

	app.Delete("/todos/:id", func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("ID"))
		if err != nil {
			return err
		}
		DB.Delete(&Todo{}, id)
		var todo []Todo
		DB.Find(&todo)
		return c.Render("todo-list", fiber.Map{
			"Todos": todo,
		})
	})

	app.Listen(":3000")
}
