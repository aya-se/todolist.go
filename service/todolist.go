package service

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	database "todolist.go/db"
	. "todolist.go/db"
	"golang.org/x/crypto/bcrypt"
	"time"
	"regexp"
	"strings"
)

var LoginInfo = User{}

// Home renders index.html
func Home(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	t := time.Now().Local().Format("2006-01-02T15:04")
	ctx.HTML(http.StatusOK, "index.html", gin.H{"Title": "HOME", "Now": t, "User": LoginInfo.UserName})
}

// TaskList renders list of tasks in DB
func TaskList(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	// タスクの検索条件指定に対応
	var searchForm SearchForm
	ctx.Bind(&searchForm)
	status := searchForm.Status
	priority := searchForm.Priority
	order := searchForm.Order
	if order == "" {
		order = "deadline"
	}
	// 指定の部分文字列
	substring := searchForm.Substring
	query_status := ""
	query_priority := ""
	// 完了状態の絞り込みクエリを追加
	if status == "incomplete" {
		query_status = " AND is_done=0"
	} else if status == "completed" {
		query_status = " AND is_done=1"
	}
	// 優先度の絞り込みクエリを追加
	if priority == "high" {
		query_priority = " AND (priority='高' OR priority='緊急')"
	}
	// Get tasks in DB
	var tasks []database.Task
	query := "SELECT * FROM tasks WHERE user_id='"+ strconv.FormatUint(LoginInfo.UserID, 10) + "' AND title LIKE '%"+ substring + "%'" + query_status + query_priority + " ORDER BY " + order
	
	err = db.Select(&tasks, query)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Render tasks
	ctx.HTML(http.StatusOK, "task_list.html", gin.H{"Title": "Task list", "Tasks": tasks, "User": LoginInfo.UserName, "Status": status, "Priority": priority, "Substring": substring, "Order": order})
}

// ShowTask renders a task with given ID
func ShowTask(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Get a task with given ID
	var task database.Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Render task
	ctx.String(http.StatusOK, task.Title)
}

//　新しいタスクの追加
func InsertTask(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// 新しいタスクの追加
	var data TaskForm
	ctx.Bind(&data)
	user_id := LoginInfo.UserID
	title := data.Title
	detail := data.Detail
	priority := data.Priority
	deadline := data.Deadline
	_, err = db.Query("INSERT INTO tasks (user_id, title, detail, priority, deadline) VALUES (?, ?, ?, ?, ?)", user_id, title, detail, priority, deadline)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
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
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
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
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
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
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Get a task with given ID
	_, err = db.Query("DELETE FROM tasks WHERE id=?", id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
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
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	status, err := strconv.Atoi(ctx.Param("status"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Get a task with given ID
	_, err = db.Query("UPDATE tasks SET is_done=? WHERE id=?", status, id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// リダイレクト
	ctx.Redirect(303, "/list")
}

// 指定タスクの編集ページ
func EditTask(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}

	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// parse ID given as a parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// Get a task with given ID
	var task Task
	err = db.Get(&task, "SELECT * FROM tasks WHERE id=?", id) // Use DB#Get for one entry
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	//data.Deadline = task.Deadline
	// Render task
	deadline := task.Deadline.Format("2006-01-02T15:04")
	ctx.HTML(http.StatusOK, "edit_task.html", gin.H{"ID": task.ID, "Title": task.Title, "Detail": task.Detail, "Priority": task.Priority, "Deadline": deadline, "User": LoginInfo.UserName })
}

// ユーザー登録ページ
func Signup(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signup.html", gin.H{"Title": "SignUp", "User": LoginInfo.UserName})
}

// ユーザーログインページ
func Signin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "signin.html", gin.H{"Title": "SignIn", "User": LoginInfo.UserName})
}

// ユーザー編集ページ
func EditUser(ctx *gin.Context) {
	//非ログイン時はリダイレクト
	if LoginInfo.UserID <= 0 {
		ctx.Redirect(303, "/signin")
		return
	}
	ctx.HTML(http.StatusOK, "edit_user.html", gin.H{"Title": "EditUser", "User": LoginInfo.UserName, "UserName": LoginInfo.UserName})
}

// ユーザー登録
func InsertUser(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// 新しいユーザーの追加
	var data UserForm
	ctx.Bind(&data)
	user_name := data.UserName
	password := data.Password
	confirm := data.Confirm

	// パスワードの一致確認
	if password != confirm {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "2つのフォームに入力されたパスワードが異なっています！再度パスワードを登録してください。"})
		return
	}
	// バリデーション
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if len(user_name)<4 || len(user_name)>16 || !re.MatchString(user_name) {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "ユーザー名には4文字以上16文字以下の英数字とアンダーバー(_)を入力してください。"})
		return
	}
	if len(password)<4 || len(password)>16 || !re.MatchString(password) {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "パスワードには4文字以上16文字以下の英数字とアンダーバー(_)を入力してください。"})
		return
	}
	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),12)
	_, err = db.Query("INSERT INTO users (user_name, password) VALUES (?, ?)", user_name, hashedPassword)
	if err != nil {
		if strings.Count(err.Error(),"Duplicate")>0 {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "このユーザー名は既に登録されています！異なるユーザー名を指定してください。"})
			return
		} else {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
			return
		}
	}

	// リダイレクト
	ctx.Redirect(303, "/signout-user")
}

// ユーザーログイン
func SigninUser(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	var data UserForm
	ctx.Bind(&data)
	user_name := data.UserName
	password := data.Password

	var user User
	err = db.Get(&user, "SELECT * FROM users WHERE user_name=?", user_name)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "ユーザー名もしくはパスワードが間違っています。"})
		return
	}
	// 退会確認
	if user.IsDeleted {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "既に退会したユーザーです。別のアカウントでログインしてください。"})
		return
	}
	// パスワード確認
	hash := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "ユーザー名もしくはパスワードが間違っています。"})
		return
	}
	// リダイレクト
	LoginInfo = user
	ctx.Redirect(303, "/list")
}

// ユーザー更新
func UpdateUser(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}

	// ユーザー情報更新
	var data UserForm
	ctx.Bind(&data)
	user_id := LoginInfo.UserID
	user_name := data.UserName
	password := data.Password
	confirm := data.Confirm

	// パスワードの一致確認
	if password != confirm {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "2つのフォームに入力されたパスワードが異なっています！再度パスワードを登録してください。"})
		return
	}
	// バリデーション
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if len(user_name)<4 || len(user_name)>16 || !re.MatchString(user_name) {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "ユーザー名には4文字以上16文字以下の英数字とアンダーバー(_)を入力してください。"})
		return
	}
	if len(password)<4 || len(password)>16 || !re.MatchString(password) {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "パスワードには4文字以上16文字以下の英数字とアンダーバー(_)を入力してください。"})
		return
	}
	// パスワードのハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password),12)
	_, err = db.Query("UPDATE users SET user_name=?, password=? WHERE user_id=?", user_name, hashedPassword, user_id)
	if err != nil {
		if strings.Count(err.Error(),"Duplicate")>0 {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": "このユーザー名は既に登録されています！異なるユーザー名を指定してください。"})
			return
		} else {
			ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
			return
		}
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
	ctx.Redirect(303, "/signin")
}

// ユーザー退会
func DeleteUser(ctx *gin.Context) {
	// Get DB connection
	db, err := database.GetConnection()
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	// ユーザーの退会処理
	user_id := LoginInfo.UserID
	_, err = db.Query("UPDATE users SET is_deleted=1 WHERE user_id=?", user_id)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"Title": "Error", "User": LoginInfo.UserName, "Error": err.Error()})
		return
	}
	// リダイレクト
	var emptyUser User
	LoginInfo = emptyUser
	ctx.Redirect(303, "/signin")
}
