package service

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	database "todolist.go/db"
	. "todolist.go/db"
	"log"
)

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT * FROM tasks ORDER BY deadline") // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Render task
	ctx.String(http.StatusOK, task.Title)
}

//　新しいタスクの追加
func InsertTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 新しいタスクの追加
	var data TaskForm
	ctx.Bind(&data)
	userid := 20000729
	title := data.Title
	detail := data.Detail
	priority := data.Priority
	deadline := data.Deadline
	_, err = db.Query("INSERT INTO tasks (user_id, title, detail, priority, deadline) VALUES (?, ?, ?, ?, ?)", userid, title, detail, priority, deadline)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/list")
}

//　指定タスクの編集
func UpdateTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	// 指定タスクの編集
	var data TaskForm
	ctx.Bind(&data)
	title := data.Title
	detail := data.Detail
	priority := data.Priority
	deadline := data.Deadline
	_, err = db.Query("UPDATE tasks SET title=? , detail=?, priority=?, deadline=? WHERE id=?", title, detail, priority, deadline, id)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/list")
}

// 指定タスクの削除
func DeleteTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get a task with given ID
	_, err = db.Query("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/list")
}

// 指定タスクの完了・再開
func CompleteTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// parse ID given as a parameter
	status, err := strconv.Atoi(ctx.Param("status"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get a task with given ID
	_, err = db.Query("UPDATE tasks SET is_done=? WHERE id=?", status, id)
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/list")
}

// 指定タスクの編集ページ
func EditTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// Get a task with given ID
	var task Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	//data.Deadline = task.Deadline
	// Render task
	deadline := task.Deadline.Format("2006-01-02T15:04")
	log.Println("%s",deadline)
	ctx.HTML(http.StatusOK, "edit_task.html", gin.H{"ID": task.ID, "Title": task.Title, "Detail": task.Detail, "Priority": task.Priority, "Deadline": deadline })
}
