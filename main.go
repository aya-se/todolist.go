package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"todolist.go/db"
	"todolist.go/service"
)

const port = 8000

func main() {
	// initialize DB connection
	dsn := db.DefaultDSN(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))
	if err := db.Connect(dsn); err != nil {
		log.Fatal(err)
	}

	// initialize Gin engine
	engine := gin.Default()
	engine.LoadHTMLGlob("views/*.html")

	// routing
	engine.Static("/assets", "./assets")
	engine.GET("/", service.Home)
	engine.GET("/list", service.TaskList)
	engine.POST("/list", service.TaskList)
	engine.GET("/task/:id", service.ShowTask)
	engine.GET("/signup", service.Signup)
	engine.GET("/signin", service.Signin)
	engine.GET("/edit-user", service.EditUser)
	engine.POST("/edit-task/:id", service.EditTask)
	engine.POST("/insert-task", service.InsertTask)
	engine.POST("/update-task/:id", service.UpdateTask)
	engine.POST("/delete-task/:id", service.DeleteTask)
	engine.POST("/complete-task/:id/:status", service.CompleteTask)
	engine.POST("/insert-user", service.InsertUser)
	engine.POST("/signin-user", service.SigninUser)
	engine.POST("/update-user", service.UpdateUser)
	engine.GET("/signout-user", service.SignoutUser)
	engine.GET("/delete-user", service.DeleteUser)
	// start server
	engine.Run(fmt.Sprintf(":%d", port))
}
