package service

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	database "todolist.go/db"
	. "todolist.go/db"
	"log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var LoginInfo = User{}

// Home renders index.html
func Home(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID == "" {
		ctx.Redirect(303, "/signin")
		return
	}
	t := time.Now().Local().Format("2006-01-02T15:04")
	ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "HOME", "Now": t, "User": LoginInfo.UserID})
}

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID == "" {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Get tasks in DB
	var tasks []database.Task
	err = db.Select(&tasks, "SELECT * FROM tasks WHERE user_id=? ORDER BY deadline", LoginInfo.UserID) // Use DB#Select for multiple entries
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks, "User": LoginInfo.UserID})
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
	//非ログイン時はリダイレクト
	if LoginInfo.UserID == "" {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 新しいタスクの追加
	var data TaskForm
	ctx.Bind(&data)
	userid := LoginInfo.UserID
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
	ctx.HTML(http.StatusOK, "edit_task.html", gin.H{"ID": task.ID, "Title": task.Title, "Detail": task.Detail, "Priority": task.Priority, "Deadline": deadline, "User": LoginInfo.UserID })
}

// ユーザー登録ページ
func Signup(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", gin.H{"Title": "SignUp", "User": LoginInfo.UserID})
}

// ユーザーログインページ
func Signin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signin.html", gin.H{"Title": "SignIn", "User": LoginInfo.UserID})
}

// ユーザー編集ページ
func EditUser(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "edit_user.html", gin.H{"Title": "EditUser", "User": LoginInfo.UserID, "UserName": LoginInfo.UserName})
}

// ユーザー登録
func InsertUser(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 新しいユーザーの追加
	var data UserForm
	ctx.Bind(&data)
	user_id := data.UserID
	user_name := data.UserName
	password, err := bcrypt.GenerateFromPassword([]byte(data.Password),12)
	_, err = db.Query("INSERT INTO users (user_id, user_name, password) VALUES (?, ?, ?)", user_id, user_name, password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/signout-user")
}

// ユーザーログイン
func SigninUser(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	var data UserForm
	ctx.Bind(&data)
	user_id := data.UserID
	password := data.Password

	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var user User
	err = db.Get(&user, "SELECT * FROM users WHERE user_id=?", user_id) // Use DB#Get for one entry
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	// パスワード確認
	hash := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		log.Println("パスワードが一致しません！")
		ctx.Redirect(303, "/signin")
		return
	}
	// リダイレクト
	LoginInfo = user
	log.Println("%v",LoginInfo)
	ctx.Redirect(303, "/list")
}

// ユーザー更新
func UpdateUser(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 新しいユーザーの追加
	var data UserForm
	ctx.Bind(&data)
	user_id := LoginInfo.UserID
	user_name := data.UserName
	_, err = db.Query("UPDATE users SET user_name=? WHERE user_id=?", user_name, user_id)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	// リダイレクト
	LoginInfo.UserName = user_name
	ctx.Redirect(303, "/list")
}

// ユーザーログアウト
func SignoutUser(ctx *gin.Context) {
	// リダイレクト
	var emptyUser User
	LoginInfo = emptyUser
	log.Println("%v",LoginInfo)
	ctx.Redirect(303, "/signin")
}
