package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Todo is a struct that represents a todo item
type Todo struct {
	ID     uint   `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

// DB is a global variable that holds the database connection
var DB *gorm.DB

func main() {
	// Connect to the database
	var err error
	DB, err = gorm.Open("sqlite3", "todo.db")
	if err != nil {
		panic("failed to connect to database")
	}
	defer DB.Close()

	// Create the todo table if it does not exist
	DB.AutoMigrate(&Todo{})

	// Create a gin router with default middleware
	r := gin.Default()

	// Define the API endpoints
	r.GET("/todos", GetTodos)          // Get all todos
	r.POST("/todos", CreateTodo)       // Create a new todo
	r.GET("/todos/:id", GetTodo)       // Get a single todo by id
	r.PUT("/todos/:id", UpdateTodo)    // Update a todo by id
	r.DELETE("/todos/:id", DeleteTodo) // Delete a todo by id

	// Start the server
	r.Run()
}

// GetTodos is a handler function that returns all todos in JSON format
func GetTodos(c *gin.Context) {
	// Create an empty slice of todos
	var todos []Todo

	// Find all todos and store them in the slice
	DB.Find(&todos)

	// Return the slice as JSON
	c.JSON(200, todos)
}

// CreateTodo is a handler function that creates a new todo from JSON input
func CreateTodo(c *gin.Context) {
	// Create an empty todo
	var todo Todo

	// Bind the JSON input to the todo
	if err := c.BindJSON(&todo); err != nil {
		// If there is an error, return a bad request status
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Save the todo to the database
	DB.Create(&todo)

	// Return the todo as JSON
	c.JSON(201, todo)
}

// GetTodo is a handler function that returns a single todo by id in JSON format
func GetTodo(c *gin.Context) {
	// Get the id parameter from the URL
	id := c.Param("id")

	// Create an empty todo
	var todo Todo

	// Find the todo by id and store it in the todo
	if err := DB.Where("id = ?", id).First(&todo).Error; err != nil {
		// If there is an error, return a not found status
		c.JSON(404, gin.H{"error": "record not found"})
		return
	}

	// Return the todo as JSON
	c.JSON(200, todo)
}

// UpdateTodo is a handler function that updates a todo by id from JSON input
func UpdateTodo(c *gin.Context) {
	// Get the id parameter from the URL
	id := c.Param("id")

	// Create an empty todo
	var todo Todo

	// Find the todo by id and store it in the todo
	if err := DB.Where("id = ?", id).First(&todo).Error; err != nil {
		// If there is an error, return a not found status
		c.JSON(404, gin.H{"error": "record not found"})
		return
	}

	// Bind the JSON input to the todo
	if err := c.BindJSON(&todo); err != nil {
		// If there is an error, return a bad request status
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Save the updated todo to the database
	DB.Save(&todo)

	// Return the todo as JSON
	c.JSON(200, todo)
}

// DeleteTodo is a handler function that deletes a todo by id
func DeleteTodo(c *gin.Context) {
	// Get the id parameter from the URL
	id := c.Param("id")

	// Create an empty todo
	var todo Todo

	// Find the todo by id and store it in the todo
	if err := DB.Where("id = ?", id).First(&todo).Error; err != nil {
		// If there is an error, return a not found status
		c.JSON(404, gin.H{"error": "record not found"})
		return
	}

	// Delete the todo from the database
	DB.Delete(&todo)

	// Return a no content status
	c.Status(204)
}
